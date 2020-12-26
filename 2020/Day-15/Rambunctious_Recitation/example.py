#!/usr/bin/env python

initial_calls = [int(num) for num in open("sample.txt").read().split(",")]
last_call = initial_calls[-1]

N  = 2020
calls = [0]*N

# Choosing indices as values and elements as turns
for index, element in enumerate(initial_calls[:-1], start=1):
  calls[element] = index

for turn in range(len(initial_calls), N):
  index = calls[last_call]
  if not index:
    index = turn
  calls[last_call], last_call = turn, turn - index

print(last_call)
