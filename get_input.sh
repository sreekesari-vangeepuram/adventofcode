#!/bin/bash

YEAR=$1
curl --cookie session=$(cat $HOME/adventofcode.com/.session-cookie) "https://adventofcode.com/$YEAR/day/$2/input"
