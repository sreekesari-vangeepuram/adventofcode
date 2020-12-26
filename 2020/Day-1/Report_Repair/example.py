#!/usr/bin/env python

numbers = [
	1721,
	979,
	366,
	299,
	675,
	1456,
]

# Identify pairs
pairs = list()
for i in numbers:
	for j in numbers:
		if i+j == 2020:
			pairs.append((i, j))
# Remove redundant pairs
for pair in pairs:
	i, j = pair
	if (j, i) in pairs:
		pairs.remove((j, i))

# Print the answer[s]
for pair in pairs:
	print(pair[0]*pair[1])
