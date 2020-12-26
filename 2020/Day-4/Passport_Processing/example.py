#!/usr/bin/env python

def verify_passport(passport):
	fields = ["byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"] # `cid` is excluded since it is not necessary
	return all(field in passport for field in fields)

input_data = open("sample.txt")

passports = list()
passport = str()
for line in input_data:
	passport += line
	if line == "\n":
		passports.append(passport.strip())
		passport = str()

print(f"The total number of valid passports are: {list(map(verify_passport, passports)).count(True)}.")

input_data.close()
