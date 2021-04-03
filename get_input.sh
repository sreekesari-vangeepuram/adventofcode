#!/bin/bash

YEAR=$1
DAY=$2
curl --cookie session=$(cat $HOME/adventofcode/.session-cookie) "https://adventofcode.com/$YEAR/day/$DAY/input"
