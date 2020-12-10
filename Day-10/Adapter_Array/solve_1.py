#!/usr/bin/env python

joltages = sorted(int(num) for num in open("input.txt"))
differences = [joltages[i+1] - joltages[i] for i in range(len(joltages)-1)]

count = lambda x: differences.count(x) + 1

print("Sample-1:", {"1": count(1), "3": count(3)})
print("Result-1:", count(1)*count(3))
