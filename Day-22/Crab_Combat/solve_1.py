#!/usr/bin/env python

from typing import List, Tuple

def play_space_cards(p1: List[int], p2: List[int]) -> Tuple[str, List[int]]:
    b1, b2 = 0, 0 # buffer spaces for both players to space their cards
    while len(p1) !=0 and len(p2)!= 0:
        b1, b2 = p1.pop(0), p2.pop(0)
        if b1 > b2:
            p1.extend([b1, b2])
        else:
            p2.extend([b2, b1])

    if len(p1) != 0:
        return "Player_1", p1
    return "Player_2", p2

def count_score(winner_deck: List[int]) -> int:
    accumulator = 0
    for card, multiplier in zip(winner_deck, list(reversed(range(1, len(winner_deck)+1)))):
        accumulator += card * multiplier
    return accumulator

decks = open("input.txt").read().strip().split("\n\n")
player_1 = list(map(int, decks[0].split("\n")[1:]))
player_2 = list(map(int, decks[1].split("\n")[1:]))

winner, winner_deck = play_space_cards(player_1, player_2)
print(f"Combat: {winner} won with score {count_score(winner_deck)}!")
