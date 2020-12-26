#!/usr/bin/env python

from collections import defaultdict

# nw, se in YZ plane
# ne, sw in XZ plane
# w ,  e in XY plane
position = {
    "nw": (0, +1, -1), "ne": (+1, 0, -1),
    "w" : (-1, +1, 0), "e" : (+1, -1, 0),
    "sw": (-1, 0, +1), "se": (0, -1, +1),
}
# `position` source: https://www.redblobgames.com/grids/hexagons/

fmt_line = lambda line: line \
                        .replace("e", "e ") \
                        .replace("w", "w ") \
                        .split()
ins_list = list(map(fmt_line, open("input.txt").read().strip().split("\n")))

tiles = defaultdict(int)
for ins in ins_list:
    x = y = z = 0
    for dx, dy, dz in [position[_in] for _in in ins]:
        x += dx; y += dy; z += dz
    tiles[x, y, z] ^= 1

print(f"Number of black sides facing-up: {sum(tiles.values())}")

