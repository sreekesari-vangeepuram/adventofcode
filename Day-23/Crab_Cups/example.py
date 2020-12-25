#!/usr/bin/env python

from itertools import count

def crab_shuffles(cups):

    # Circle is a dictionary [used in place of linked_list]
    circle = {cup: next_cup for cup, next_cup in zip(cups, [*cups[1:], cups[0]])}
    current_cup = cups[0]

    for _ in range(MOVES):
        cp = current_cup # cp -> cup-pointer
        pick_up = [cp := circle[cp] for _ in range(3)] # traversing the next 3 nodes
        destination = next(cup for i in count(1) if (cup if int(cup := current_cup - i) > 0 else int(cup := len(cups) + cup)) not in pick_up)

        (circle[current_cup],
         circle[pick_up[-1]],
         circle[destination]) = (circle[pick_up[-1]],
                                 circle[destination],
                                 circle[current_cup])

        current_cup = circle[current_cup]

    cp = 1
    return "".join(map(str, [cp := circle[cp] for _ in cups][:-1]))

cups_list = list(map(int, open("sample.txt").read().strip()))
global MOVES
MOVES = 100

print(f"Labels on the cup after 100th move: {crab_shuffles(cups_list)}")

