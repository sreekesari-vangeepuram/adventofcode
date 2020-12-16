#!/usr/bin/env python

import re

input_data = open("input.txt").read().strip().split("\n\n")

# Extracting required data from the <List> -> `input_data`
flight_info = dict(tuple(line.split(": ")) for line in input_data[0].split("\n"))
my_ticket = tuple(int(num) for num in input_data[1].strip().split("\n")[1].split(","))
tickets = [my_ticket] + [tuple(int(num) for num in line.split(",")) for line in input_data[2].strip().split("\n")[1:]]

# Parsing and assigning of valid range set to the given fields
for key, value in flight_info.items():
  buffer = list(int(num) for num in re.findall(r"(\d+)-(\d+) or (\d+)-(\d+)", value)[0])
  flight_info[key] = set(range(buffer[0], buffer[1] + 1)) | set(range(buffer[2], buffer[3] + 1))

# Validation of tickets
valid_tickets = set()
for ticket in tickets:
  valid = True
  for number in ticket:
    if all(number not in flight_info[field] for field in flight_info):
      valid = False
  if valid:
    valid_tickets.add(ticket)


# Determining mappings
# Pairing fields of the dictionary `flight_info` in correct order
buffer_map = {field: set() for field in flight_info.keys()}
for field in buffer_map.keys():
  for index, field_column in enumerate(zip(*valid_tickets)):
    if all(True if number in flight_info[field] else False for number in field_column):
      buffer_map[field].add(index)

# Also can be done in this way
"""
index = 0
while index < len(buffer_map):
  possible_fields = set(buffer_map.keys())
  for ticket in valid_tickets:
    possible_fields = set(field for field in possible_fields if ticket[index] in flight_info[field])
  for field in possible_fields:
    buffer_map[field].add(index)
  index += 1
"""

pairs = {}
while buffer_map:
  for field in buffer_map:
    if len(buffer_map[field]) == 1: # 1st iteration: <class> field of `buffer_map`
      pairs[field] = buffer_map[field].pop()
      buffer_map.pop(field)
      # Removal of determined field_column from the `buffer_map`
      for f in buffer_map:
        buffer_map[f].remove(pairs[field])
      break

accumulator = 1
for field_name, index in pairs.items():
  if "departure" in field_name:
    accumulator *= my_ticket[index]

print(f"Product of values with `departure` in their field name: {accumulator}")
