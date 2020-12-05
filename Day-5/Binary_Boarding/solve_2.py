#!/usr/bin/env python
"""
def getX(row, low, high, terms):
	mid = low + (high - low)//2
	for i, x in enumerate(row):
		if x == terms[0]: # F or L
			return getX(row[i+1:], low, mid-1, terms)
		elif x == terms[1]: # B or R
			return getX(row[i+1:], mid+1, high, terms)
	return mid+1
"""

input_data = open("input.txt")

"""
seats = list()
for seat in input_data:
	id = getX(seat[:7], 0, 127, ("F", "B"))*8 + getX(seat[7:], 0, 7, ("L", "R"))
	seats.append(id)

input_data.seek(0)
"""
bit_map = {"B": 1, "F": 0, "R": 1, "L": 0} # Since `F` is lower-half and `B` is upper-half and same for `L`, `R` respectively!

seats = list()
for string in input_data:
	seats.append(int("".join([str(bit_map[letter]) for letter in string.strip()]), 2))

unfilled_seats = list(seat for seat in range(min(seats), max(seats)+1) if seat not in seats)
print(unfilled_seats)

input_data.close()
