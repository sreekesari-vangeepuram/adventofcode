#!/usr/bin/env python

input_data = list(map(int, open("sample.txt").read().strip().split("\n")))

for i, num in enumerate(input_data[5:]):
	combinations = set(i1+i2 for i1 in input_data[i:i+5] for i2 in input_data[i:i+5])
	if num not in combinations:
		print(num)
