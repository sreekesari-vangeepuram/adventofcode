#!/usr/bin/env python

input_data = open("input.txt")

bit_map = {"B": 1, "F": 0, "R": 1, "L": 0} # Since `F` is lower-half and `B` is upper-half and same for `L`, `R` respectively!

seats = list(int("".join([str(bit_map[letter]) for letter in string.strip()]), 2) for string in input_data)
unfilled_seats = list(seat for seat in range(min(seats), max(seats)+1) if seat not in seats)

print(unfilled_seats)

input_data.close()
