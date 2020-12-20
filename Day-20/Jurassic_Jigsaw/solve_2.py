#!/usr/bin/env python

def assemble_tiles(tiles, corner_ids):

    for corner_id in corner_ids:

        buffer_dict = dict(tiles)
        _id = buffer_dict.pop(corner_id)

        l_side = map_rows(buffer_dict,
                    _id,
                    lambda prev, pres: prev[-1] == pres[0])

        rows = list(map_rows(buffer_dict,
                    tile,
                    lambda prev, pres: all(row_1[-1] == row_2[0] for row_1, row_2 in zip(prev, pres))) for tile in l_side)

        if not buffer_dict:
            return list("".join(line[1:-1] for line in chunk)
                        for row in rows for chunk in zip(*(unit[1:-1] for unit in row)))

    Exception("Unable assemble the tiles!")


def map_rows(buff_dict, tile_0, fn):
    mapped_rows = [tile_0]
    while True:
        mapped = False
        for _id, tile in tuple(buff_dict.items()):
            for flipped_tile in flip(tile):
                if fn(mapped_rows[-1], flipped_tile):
                    mapped_rows.append(flipped_tile)
                    mapped = True
                    _ = buff_dict.pop(_id)
                    break

        if not mapped:
            return mapped_rows

    Exception("Unable to map rows!")


def get_sides(tile):
    sides = [
        tile[0],                                        # Top
        tile[-1],                                       # Bottom
        "".join(list(row[0] for row in tile)),          # l-side
        "".join(list(row[0] for row in tile[::-1])),    # r-side
    ]

    T_sides = list(row[::-1] for row in sides)          # Transpose of `sides`

    return list(sides + T_sides)

def get_sides_of_all(tiles):
    buffer_dict = {key: set() for key in tiles.keys()}
    for id, tile in tiles.items():
        buffer_dict[id].add(tile[0])         # top
        buffer_dict[id].add(tile[-1])        # bottom
        tile_sides = ["".join(t) for t in list(zip(*tile))]
        buffer_dict[id].add(tile_sides[0])   # l-side
        buffer_dict[id].add(tile_sides[-1])  # r-side
    return buffer_dict

def flip(tile):
    rotate_4 = 0
    # Since each tile can rotate in 8 times
    while rotate_4 < 2:
        yield tile                                                     #0 no-change in I-1, columns -> rows in I-2
        yield tile[::-1]                                               #1 flip #0 over X-axis
        yield list(row[::-1] for row in tile)                          #2 flip #0 over Y-axis
        yield list(row[::-1] for row in tile[::-1])                    #3 flip #2 over X-axis
        tile = list("".join(column) for column in zip(*tile))          #4 Transpose [rows   ->  columns]
        rotate_4 += 1

tiles = [tile.split("\n") for tile in open("input.txt").read().strip().split("\n\n")]
tiles = {int(tile[0][5:-1]): tile[1:] for tile in tiles}

sides_set_of_tiles = get_sides_of_all(tiles) # {id: {top, bottom, l-side, r-side}}
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
least_count = min(match_count.values())
corner_ids = list()
for tile_id, count in match_count.items():
    if count == least_count:
        corner_ids.append(tile_id)


image = assemble_tiles(tiles, corner_ids)

monster_image = [
    "                  # ",
    "#    ##    ##    ###",
    " #  #  #  #  #  #   ",
]

image_len = len(image)
image_row_len = len(image[0])

monster_parts = set()
for flipped_monster in flip(monster_image):
    monster_image_len = len(flipped_monster)
    monster_image_row_len = len(flipped_monster[0])

    active_region = list()
    for dx, row in enumerate(flipped_monster):
        for dy, column in enumerate(row):
            if column == "#":
                active_region += [(dx, dy)]

    for x in range(image_len - monster_image_len):
        for y in range(image_row_len - monster_image_row_len):

            parts = {
                (x + dx, y + dy) for (dx, dy) in active_region
            }

            if all(image[i][j] == "#" for i, j in parts):
                monster_parts |= parts                              # Extracting the monster `parts` from `active_region`


accumulator = 0
for row in image:
    for column in row:
        if column == "#":
            accumulator += 1
# Since, we have to submit the number of hash `#`
# symbols which are not part of the sea monster
hash_count = accumulator - len(monster_parts)

print(f"`#` count without the parts of sea monster: {hash_count}")
