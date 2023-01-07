package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/tarm/serial"
)

var port = flag.String("port", "/dev/ttyUSB0", "Serial port")
var baud = flag.Int("baud", 115200, "Baud rate")
var up = flag.Bool("up", false, "Just raise pen")
var down = flag.Bool("down", false, "Just lower pen")

func main() {
	flag.Parse()

	if *up {
		cmdAndExit("m5")
	}

	if *down {
		cmdAndExit("m3 s90")
	}

	args := flag.Args()
	if len(args) < 1 {
		log.Printf("usage: %s <file.grbl>", os.Args[0])
		flag.Usage()
	}

	cmdFile, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}
	defer cmdFile.Close()

	cmdR := bufio.NewReader(cmdFile)

	conf := &serial.Config{Name: *port, Baud: *baud}
	s, err := serial.OpenPort(conf)
	if err != nil {
		panic(err)
	}

	fromDevice := bufio.NewReader(s)
	toDevice := bufio.NewWriter(s)

	toDevice.WriteString("~\n")
	toDevice.Flush()

	log.Printf("startup loop")
	for {
		log.Printf("read line")
		c, err := fromDevice.ReadBytes('\n')
		if err != nil {
			log.Fatalf("Read from device err: %s", err)
		}
		log.Printf("got line: %s", c)
		m := string(c)
		if len(m) == 26 && m[:5] == "Grbl " && m[9:] == " ['$' for help]\r\n" {
			fmt.Printf("Grbl version %s initialized\n", m[5:9])
			break
		} else if m == "ok\r\n" {
			break
		} else if m == "\r\n" {
			continue
		} else {
			log.Printf("got %s\n", c)
		}
	}

	for {
		line, err := cmdR.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "(") {
			log.Printf("#: <%s>", line)
			continue
		}

		log.Printf("send: <%s>", line)

		_, err = toDevice.WriteString(line + "\n")
		if err != nil {
			log.Fatalf("Write err: %s", err)
		}
		err = toDevice.Flush()
		if err != nil {
			log.Fatalf("Write err: %s", err)
		}

		resp, err := fromDevice.ReadString('\n')
		if err != nil {
			log.Fatalf("read response err: %s", err)
		}
		if resp == "ok\r\n" {
			log.Printf("ok")
			// ok!
		} else if len(resp) >= 5 && resp[:5] == "error" {
			log.Printf("got err from device: %s", resp)
		} else if len(resp) >= 5 && resp[:5] == "alarm" {
			log.Printf("got alarm from device: %s", resp)
		} else {
			log.Printf("got other from device: %s", resp)
		}
	}

}

func cmdAndExit(cmd string) {
	conf := &serial.Config{Name: *port, Baud: *baud}
	s, err := serial.OpenPort(conf)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(s, "~\n")
	fmt.Fprintf(s, "%s\n", cmd)

	buf := make([]byte, 128)
	n, _ := s.Read(buf)

	fmt.Printf("%s\n", buf[:n])
	os.Exit(0)
}
