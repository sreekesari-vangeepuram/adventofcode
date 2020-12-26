#!/usr/bin/env python

import re

BAG = "shiny gold"
input_data = open("input.txt")

bag_container = dict() # dictionaries
for line in input_data:
	tmp = line.strip().split("contain") # temporary buffer variable for the splitted line
	bag_container[tmp[0].replace("bags", "").strip()] = list(map(lambda x: x.replace("bags", "").replace("bag", "").replace(".", "").strip(), tmp[1].split(",")))

bag_container = {k: {x[2:].strip(): int(x[:2].strip()) for x in v if x != "no other"} for k, v in bag_container.items()}

get_lim = lambda container, colour: sum(v + v * get_lim(container, k) for k, v in bag_container[colour].items())
print(get_lim(bag_container, BAG))

input_data.close()
