#!/usr/bin/env python

from itertools import count

def crab_shuffles(cups):

    # Circle is a dictionary [used in place of linked_list]
    circle = {cup: next_cup for cup, next_cup in zip(cups, [*cups[1:], cups[0]])}
    current_cup = cups[0]

    for _ in range(MOVES):
        cp = current_cup # cp -> cup-pointer
        pick_up = [cp := circle[cp] for _ in range(3)] # traversing the next 3 nodes
        destination = next(cup for i in count(1) if (cup if (cup := current_cup - i) > 0 else (cup := len(cups) + cup)) not in pick_up)

        (circle[current_cup],
         circle[pick_up[-1]],
         circle[destination]) = (circle[pick_up[-1]],
                                 circle[destination],
                                 circle[current_cup])

        current_cup = circle[current_cup]

    cp = 1
    return [cp := circle[cp] for _ in cups]

cups_list = list(map(int, open("input.txt").read().strip()))
global MOVES
MOVES = 10 * (10**6) # 10 million moves

buffer_list = cups_list + list(range(len(cups_list) + 1, 10**6 + 1)) # +1M Cups

result = crab_shuffles(buffer_list)
print(f"Labels on the cup after 10Mth move: {result[0] * result[1]}")

