#!/usr/bin/env python

from functools import reduce

joltages = sorted(int(num) for num in open("input.txt"))
sequence_split = "".join([str(max(1, joltage-joltages[i-1])) for i, joltage in enumerate(joltages)]).split("3")

permutations = list()
for ones in sequence_split:
	exp = len(ones) - 1
	if exp < 3: permutations.append(max(1, 2**(exp)))
	else: permutations.append(2**(exp) - 1)

print("Count of arrangements:", reduce(lambda x, y: x*y, permutations))
