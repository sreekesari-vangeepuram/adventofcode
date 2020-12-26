#!/usr/bin/env python

import re
from functools import reduce
from collections import defaultdict

BAG = "shiny gold"
input_data = open("sample.txt")

bag_container = dict() # dictionaries
for line in input_data:
	tmp = line.strip().split("contain") # temporary buffer variable for the splitted line
	bag_container[tmp[0].replace("bags", "").strip()] = list(map(lambda x: x.replace("bags", "").replace("bag", "").replace(".", "").strip(), tmp[1].split(",")))

bag_container = {k: {x[2:].strip(): int(x[:2].strip()) for x in v if x != "no other"} for k, v in bag_container.items()}
new_bag_container = defaultdict(set) # dictionary of `string: key` and `set: value` pairs
for parent_bag, child_bags in bag_container.items():
	for child_bag in child_bags:
		new_bag_container[child_bag].add(parent_bag)

atleast_one = lambda container, colour: reduce(set.union, (atleast_one(container, c) for c in container[colour]), container[colour])
print(len(atleast_one(new_bag_container, BAG)))

input_data.close()
