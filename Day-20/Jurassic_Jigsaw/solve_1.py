#!/usr/bin/env python

def get_sides(tiles):
    buffer_dict = {key: set() for key in tiles.keys()}
    for id, tile in tiles.items():
        buffer_dict[id].add(tile[0])         # top
        buffer_dict[id].add(tile[-1])        # bottom
        tile_sides = ["".join(t) for t in list(zip(*tile))]
        buffer_dict[id].add(tile_sides[0])   # l-side
        buffer_dict[id].add(tile_sides[-1])  # r-side
    return buffer_dict


tiles = [tile.split("\n") for tile in open("input.txt").read().strip().split("\n\n")]
tiles = {int(tile[0][5:-1]): tile[1:] for tile in tiles}

sides_set_of_tiles = get_sides(tiles) # {id: {top, bottom, l-side, r-side}}
match_count = {key: 0 for key in tiles.keys()}
for id1, t1 in sides_set_of_tiles.items():
    for id2, t2 in sides_set_of_tiles.items():
        count = len(t1.intersection(t2))
        t1 = {side[::-1] for side in t1} # flipped_sides
        flip_count = len(t1.intersection(t2))
        match_count[id1] += count + flip_count
        match_count[id2] += count + flip_count

# Since corner maps with only 2 adjacent tiles
# they have least count of matches.
accumulator, least_count = 1, min(match_count.values())
for tile_id, count in match_count.items():
    if count == least_count:
        accumulator *= tile_id
print(f"Product of corner-ids: {accumulator}")
