#!/usr/bin/env python

input_data = open("sample.txt")

field = list(map(lambda x: x.strip(), input_data))
max_len = len(field[0])
field = [row*(max_len+1) for row in field]

y = 0
counter = 0
for x in range(len(field)):
	y += 3
	if x < len(field) - 1 and field[x+1][y] == "#":
		counter += 1
print(f"Encountered trees count: {counter}.")

input_data.close()
