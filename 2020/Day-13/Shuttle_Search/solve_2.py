#!/usr/bin/env python

input_data = open("input.txt").read().strip().split("\n")
bus_ids = tuple((int(id), (int(id) - i) % int(id)) for i, id in enumerate(input_data[1].split(",")) if id != 'x')

tmp = bus_ids[0][0]
timestamp, i = tmp, 0

# Chinese Remainder Theorem

loop = True
while loop:
  id, remainder = bus_ids[i+1]
  if (timestamp) % id  == remainder:
    if i == len(bus_ids) - 2:
      print(f"Timestamp of the earliest bus with subsequent offset: {timestamp}")
      loop = False
    tmp *= id
    i += 1
  timestamp += tmp
