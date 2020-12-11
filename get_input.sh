#!/bin/bash

YEAR=2020
curl --cookie session=$(cat .session-cookie) "https://adventofcode.com/$YEAR/day/$1/input"
