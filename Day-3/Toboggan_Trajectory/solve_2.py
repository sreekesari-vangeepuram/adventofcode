#!/usr/bin/env python

from functools import reduce

input_data = open("input.txt")

field = list(map(lambda x: x.strip(), input_data))
field_height, field_breadth = len(field), len(field[0])
field = [row *(73) for row in field]

slopes = ((1, 1), (3, 1), (5, 1), (7, 1), (1, 2))
results = list()

for slope in slopes:
	right, down = slope
	x, y = 0, 0
	counter = 0
	while x < field_height:
		if field[x][y] == "#":
			counter += 1
		x += down
		y += right
	results.append(counter)
print(f"Encountered trees count: {reduce(lambda x, y: x*y, results)}.")

input_data.close()
