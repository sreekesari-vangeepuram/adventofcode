#!/usr/bin/env python

input_data = open("input.txt")

field = list(map(lambda x: x.strip(), input_data))
max_len = len(field[0])
field = [row*(max_len+1) for row in field]

right, down = 3, 1
y = 0
counter = 0
for x in range(len(field)):
	y += right
	if x < len(field) - 1 and field[x+down][y] == "#":
		counter += 1
print(f"Encountered trees count: {counter}.")

input_data.close()
