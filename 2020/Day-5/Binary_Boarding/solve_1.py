#!/usr/bin/env python

def getX(row, low, high, terms):
	mid = low + (high - low)//2
	for i, x in enumerate(row):
		if x == terms[0]: # F or L
			return getX(row[i+1:], low, mid-1, terms)
		elif x == terms[1]: # B or R
			return getX(row[i+1:], mid+1, high, terms)
	return mid+1

input_data = open("input.txt")

highest_id = 0
for seat in input_data:
	id = getX(seat[:7], 0, 127, ("F", "B"))*8 + getX(seat[7:], 0, 7, ("L", "R"))
	if highest_id < id:
		highest_id = id

print(highest_id)

input_data.close()
