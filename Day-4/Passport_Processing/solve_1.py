#!/usr/bin/env python

def verify_passport(passport):
	fields = ["byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"] # `cid` is excluded since it is optional
	return all(field in passport for field in fields)

result = list(map(verify_passport, open('input.txt').read().strip().split("\n\n"))).count(True)

print(f"The total number of valid passports are: {result}.")
