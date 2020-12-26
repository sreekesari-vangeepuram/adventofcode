#!/usr/bin/env python

from itertools import product
from collections import defaultdict

def neighbours(point):
    for combination in product([-1, 0, 1], repeat=len(point)):
        neighbour = tuple(dp + combination[i] for i, dp in enumerate(point))
        yield neighbour

grid = open("input.txt").read().strip().split("\n")
DIMENSIONS = 3
ACTIVE, INACTIVE = "#", "."

z = (0,) * (DIMENSIONS - 2) # Buffer for the 3rd axis
#space = dict(((x, y, *z), cube) for x, row in enumerate(grid) for y, cube in enumerate(row))
# ^^^ This list-comprehension doesn't work since as mentioned in the challenge statement:
# All the cube have their initial-state as `INACTIVE`.
space = defaultdict(lambda: INACTIVE)
for x, row in enumerate(grid):
    for y, cube in enumerate(row):
        space[(x, y, *z)] = cube

layer_cycle = 0
while layer_cycle < 6:
    active_cubes = defaultdict(int)
    for cube in space:
        if space[cube] == INACTIVE: continue
        for neighbour in neighbours(cube):
            # <variable> += True  [increments  as int(True)]
            # <variable> += False [neutralizes as int(False)]
            active_cubes[neighbour] += neighbour != cube and space[cube] == ACTIVE

    for cube, active_count in active_cubes.items():
        if space[cube] == ACTIVE:
            space[cube] = ACTIVE if active_count in {2, 3} else INACTIVE

        elif space[cube] == INACTIVE and active_count == 3:
            space[cube] = ACTIVE

    layer_cycle += 1

print(f"Active cubes left in the 6th cycle: {sum(cube == ACTIVE for cube in space.values())}")
