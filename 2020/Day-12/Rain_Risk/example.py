#!/usr/bin/env python

from vector import vector

DIRECTIONS = ('E', 'W', 'N', 'S', 'F')
TURNINGS = ('R', 'L')

nav_ins = list((ins[0], int(ins[1:])) for ins in open('sample.txt'))
ship = vector(0, 0, 'E')

print(ship.get_pos())
for ins in nav_ins:
	if ins[0] in DIRECTIONS:
		ship.change_position(ins)
	elif ins[0] in TURNINGS:
		ship.change_direction(ins)
	print(ship.get_pos())

print(f"Manhattan distance: {ship.manhattan_distance()}")
