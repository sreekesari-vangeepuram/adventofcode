#!/usr/bin/env python

joltages_1 = sorted(int(num) for num in open("sample.txt"))
joltages_2 = sorted(int(num) for num in open("sample2.txt"))

differences_1 = [joltages_1[i+1] - joltages_1[i] for i in range(len(joltages_1)-1)]
differences_2 = [joltages_2[i+1] - joltages_2[i] for i in range(len(joltages_2)-1)]

count_1 = lambda x: differences_1.count(x) + 1
count_2 = lambda x: differences_2.count(x) + 1

print("Sample-1:", {"1": count_1(1), "3": count_1(3)})
print("Result-1:", count_1(1)*count_1(3))
print()
print("Sample-2:", {"1": count_2(1), "3": count_2(3)})
print("Result-2:", count_2(1)*count_2(3))
