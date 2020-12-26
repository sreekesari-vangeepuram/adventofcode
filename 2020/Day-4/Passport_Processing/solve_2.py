#!/usr/bin/env python

import re

cases = {
    'byr': lambda x: re.match('^\d{4}$', x) and 1920 <= int(x) <= 2002,
    'iyr': lambda x: re.match('^\d{4}$', x) and 2010 <= int(x) <= 2020,
    'eyr': lambda x: re.match('^\d{4}$', x) and 2020 <= int(x) <= 2030,
    'hgt': lambda x: (re.match('^\d{3}cm$', x) and 150 <= int(x.replace('cm', '')) <= 193) or (re.match('^\d{2}in$', x) and 59 <= int(x.replace('in', '')) <= 76),
    'hcl': lambda x: re.match('^\#[0-9a-f]{6}$', x),
    'ecl': lambda x: x in ('amb', 'blu', 'brn', 'gry', 'grn', 'hzl', 'oth'),
    'pid': lambda x: re.match('^\d{9}$', x)
}

passports = list(dict(x) for x in map(lambda li: [tuple(s.split(':')) for s in li], map(lambda x: x.split(), open('input.txt').read().strip().split("\n\n"))))

counter = 0
for passport in passports:
	counter += all(fn(passport.get(key, "")) for key, fn in cases.items())
print(f"The total number of valid passports are: {counter}.")
