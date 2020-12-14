#!/usr/bin/env python

from itertools import product

def permutated_address(value, count):
  permutations = list(product([0, 1], repeat=count))
  accumulator = 0
  for permutation in permutations:
    yield int(value.replace("X", "{}").format(*permutation), 2)

def mask_on_all(pair):
  for mask, instructions in pair.items():
    buff_dict = dict()
    for address, value in instructions.items():
      tmp = "".join([bit if mask[i] == "0" else mask[i] for i, bit in enumerate(bin(int(address[4:-1])).rjust(36).replace(" ", "0").replace("b", "0"))])
      for new_address in permutated_address(tmp, tmp.count("X")):
        buff_dict[f"mem[{new_address}]"] = int(value)
  return buff_dict

pairs = list()
for line in open("input.txt").read().strip().split("mask = "):
  buff_list = line.strip().split('\n')
  pairs.append({buff_list[0]: dict([tuple(pair.split(" = ")) for pair in buff_list[1:]])})

mem_dict = dict()
for pair in pairs:
  mem_dict.update(mask_on_all(pair))

print(f"Sum of all values left in memory: {sum(mem_dict.values())}")
