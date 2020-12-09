#!/usr/bin/env python

input_data = list(map(int, open("input.txt").read().strip().split("\n")))

for i, num in enumerate(input_data[25:]):
	combinations = set(i1+i2 for i1 in input_data[i:i+25] for i2 in input_data[i:i+25])
	if num not in combinations:
		print("Invalid num:", num)
