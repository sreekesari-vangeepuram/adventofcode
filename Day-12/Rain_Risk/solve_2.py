#!/usr/bin/env python

nav_ins = list((ins[0], int(ins[1:])) for ins in open('input.txt'))

ship = complex(0, 0)
waypoint = complex(10, 1)

direction = {
	'N': complex(0, 1),
	'S': complex(0, -1),
	'E': complex(1, 0),
	'W': complex(-1, 0)
}

for (ins, val) in nav_ins:

	# Translation of waypoint according to instructions
	if ins == "N": waypoint += val*direction["N"]
	elif ins == "S": waypoint += val*direction["S"]
	elif ins == "E": waypoint += val*direction["E"]
	elif ins == "W": waypoint += val*direction["W"]

	# Rotation of waypoint around the ship
	elif ins == "R": waypoint *= direction["S"] ** (val // 90)
	elif ins == "L": waypoint *= direction["N"] ** (val // 90)

	# Scaled translation of ship position towards the waypoint
	elif ins == "F": ship += waypoint * val

print(f"Manhattan Distance: {int(abs(ship.real)+abs(ship.imag))}")
