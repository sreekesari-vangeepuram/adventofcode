#!/usr/bin/env python

def mask_on_all(pair):
  for mask, instructions in pair.items():
    buff_dict = dict()
    for address, value in instructions.items():
      buff_dict[address] = int("".join([bit if mask[i] == "X" else mask[i] for i, bit in enumerate(bin(int(value)).rjust(36).replace(" ", "0").replace("b", "0"))]), 2)
  return buff_dict

pairs = list()
for line in open("input.txt").read().strip().split("mask = "):
  buff_list = line.strip().split('\n')
  pairs.append({buff_list[0]: dict([tuple(pair.split(" = ")) for pair in buff_list[1:]])})

mem_dict = dict()
for pair in pairs:
  mem_dict.update(mask_on_all(pair))

print(f"Sum of all values left in memory: {sum(mem_dict.values())}")
