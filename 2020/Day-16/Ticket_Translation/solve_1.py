#!/usr/bin/env python

import re

input_data = open("input.txt").read().strip().split("\n\n")

flight_info = dict(tuple(line.split(": ")) for line in input_data[0].split("\n"))
my_ticket = tuple(int(num) for num in input_data[1].strip().split("\n")[1].split(","))
tickets = tuple(int(num) for num in ",".join(input_data[2].split("\n")[1:]).strip().split(","))

field_set = set()
for value in flight_info.values():
  buffer = list(int(num) for num in re.findall(r"(\d+)-(\d+) or (\d+)-(\d+)", value)[0])
  field_set |= set(range(buffer[0], buffer[1] + 1)) | set(range(buffer[2], buffer[3] + 1))

accumulator = 0
for ticket in tickets:
  if ticket not in field_set:
    accumulator += ticket

print(f"Ticket scanning error rate: {accumulator}")
