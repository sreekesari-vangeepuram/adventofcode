#!/usr/bin/env python

from typing import List, Tuple

def play_space_cards(p1: List[int], p2: List[int]) -> Tuple[str, List[int]]:
    bs1, bs2 = set(), set() # buffer sets for storing decks
    b1, b2 = 0, 0           # buffer spaces for both players to space their cards
    while len(p1) !=0 and len(p2)!= 0:
        if tuple(p1) in bs1 or tuple(p2) in bs2:
            return "Player_1", p1
        bs1.add(tuple(p1))
        bs2.add(tuple(p2))
        b1, b2 = p1.pop(0), p2.pop(0)

        # According to new rules
        if b1 <= len(p1) and b2 <= len(p2):
            winner, _ = play_space_cards(p1[:b1], p2[:b2])
            if winner == "Player_1":
                p1.extend([b1, b2])
            elif winner == "Player_2":
                p2.extend([b2, b1])

        # According to old rules
        else:
            if b1 > b2:
                p1.extend([b1, b2])
            else:
                p2.extend([b2, b1])

    if len(p1) != 0:
        return "Player_1", p1
    return "Player_2", p2

def count_score(winner_deck: List[int]) -> int:
    accumulator = 0
    for card, multiplier in zip(winner_deck, range(len(winner_deck), 0, -1)):
        accumulator += card * multiplier
    return accumulator

decks = open("input.txt").read().strip().split("\n\n")
player_1 = list(map(int, decks[0].split("\n")[1:]))
player_2 = list(map(int, decks[1].split("\n")[1:]))

winner, winner_deck = play_space_cards(player_1, player_2)
print(f"Recursive Combat: {winner} won with score {count_score(winner_deck)}!")
