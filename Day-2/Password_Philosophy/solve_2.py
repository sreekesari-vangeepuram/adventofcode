#!/usr/bin/env python

import re

input_data = open("input.txt")
pattern = r"([0-9]{1,})-([0-9]{1,}) ([a-z]): ([a-z]{1,})"

counter = 0
for line in input_data:
	ll, ul, letter, password = re.findall(pattern, line)[0]
	if bool(password[int(ll)-1] == letter) ^ bool(password[int(ul)-1] == letter):
		counter += 1
input_data.close()

print(f"The eligible passwords count is: {counter}.")
