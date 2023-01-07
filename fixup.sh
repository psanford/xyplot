#!/usr/bin/env bash


# 's/G00 Z5.000000/M5 S90/g'          down pen
# 's/G01 Z-0.125000 F100.0/M3 s90/g'  up pen
# 's/F400.000000/F1500/g'             movement speed
# 's/G00/G01/g'                       move at normal speed for everything

sed -re 's/G00 Z5.000000/M5 S90/g' $1 | sed -re 's/G01 Z-0.125000 F100.0/M3 s90/g' | sed -re 's/F400.000000/F1500/g' | sed -re 's/G00/G01/g'
