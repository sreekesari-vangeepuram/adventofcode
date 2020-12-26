#!/usr/bin/env python

import re

input_data = [
	"1-3 a: abcde",
	"1-3 b: cdefg",
	"2-9 c: ccccccccc",
]

pattern = r"([0-9]{1,})-([0-9]{1,}) ([a-z]): ([a-z]{1,})"

counter = 0
for line in input_data:
	ll, ul, letter, password = re.findall(pattern, line)[0]
	if int(ll) <= password.count(letter) <= int(ul):
		counter += 1

print(f"The eligible passwords count is: {counter}.")
