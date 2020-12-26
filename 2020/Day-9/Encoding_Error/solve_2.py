#!/usr/bin/env python

input_data = list(map(int, open("input.txt").read().strip().split("\n")))

idx, inv_num = 0, 0
for i, num in enumerate(input_data[25:]):
	combinations = set(i1+i2 for i1 in input_data[i:i+25] for i2 in input_data[i:i+25])
	if num not in combinations:
		idx, inv_num = i, num

def get_pair(t, num):
	for i in range(t, idx):
		pairs = input_data[i-t:i]
		if sum(pairs) == num:
			return ((i-t, i), pairs)
	return (None, None)

for i in range(2, 100):
	chunk_range, pair = get_pair(i, inv_num)
	if chunk_range != None and pair != None:
		print("Encryption Weakness:", min(pair) + max(pair))
