#!/usr/bin/env python

input_data = open("sample.txt").read().strip().split("\n")
timestamp = int(input_data[0])
bus_ids = [int(id) for id in input_data[1].split(",") if id != 'x']

dep_time, bus_id = 0, 0
i = 0
while bus_id == 0:
  for id in bus_ids:
    if (timestamp + i) % id == 0:
      dep_time = timestamp + i
      bus_id = id
  i += 1
print(f"ID of the earliest bus you can take to the airport multiplied by the number of minutes: {(dep_time-timestamp)*bus_id}")
