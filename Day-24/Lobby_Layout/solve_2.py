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

def game_of_life(neighbours):
    black_tiles = defaultdict(bool)
    for cell in neighbours.keys():
        case_1 =  tiles[cell] and 0 < neighbours[cell] < 3
        case_2 = not tiles[cell] and neighbours[cell] == 2
        black_tiles[cell] = case_1 or case_2
    return black_tiles

fmt_line = lambda line: line \
                        .replace("e", "e ") \
                        .replace("w", "w ") \
                        .split()
ins_list = list(map(fmt_line, open("input.txt").read().strip().split("\n")))

global tiles # For accessing it in `game_of_life`
tiles = defaultdict(int)
for ins in ins_list:
    x = y = z = 0
    for dx, dy, dz in [position[_in] for _in in ins]:
        x += dx; y += dy; z += dz
    tiles[x, y, z] ^= 1

DAYS = 100
for _ in range(DAYS):
    neighbours = defaultdict(int) # buffer dictionary
    for x, y, z in tiles.keys():
        if tiles[x, y, z]: # (cell) is alive, then
            for dx, dy, dz in position.values():
                neighbours[x + dx, y + dy, z + dz] += 1

    tiles = game_of_life(neighbours) # update the state of tiles

print(f"Black sides facing-up on Day-100: {sum(tiles.values())}")

