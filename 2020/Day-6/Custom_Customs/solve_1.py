#!/usr/bin/env python

input_data = open("input.txt").read().split("\n\n")

groups = [set(x.strip().split("\n")) for x in input_data]

filtered_responses = list() # list of sets
tmp_set = set() # temporary set as buffer-space for each group

for group in groups: # set
	for person in group: # string
		for answer in person: # letter
			tmp_set.add(answer)
	filtered_responses.append(tmp_set)
	tmp_set = set() # killing the buffer

counter = 0
for grp_resps in filtered_responses:
	counter += len(grp_resps)
print(f"The final sum of the counts from the the input-set is: {counter}.")
