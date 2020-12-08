#!/usr/bin/env python

input_data = list(map(lambda x: x.split(), open("sample.txt").read().strip().split("\n")))

ip_1 = 0 # 1st operation record
for t in input_data:
	if t[0] != "nop":
		ip_1 = input_data.index(t)
		break

acc = 0 # accumulator holding buffer variable
ip = 0 # instruction pointer
cnt = 0 # counter to track the `ip` over `ip_1`

while True:
	if input_data[ip][0] == "acc":
		acc += int(input_data[ip][1])
	if input_data[ip][0] == "jmp":
		ip += int(input_data[ip][1]) - 1 # `-1` to balance the ip increment in last step
	if ip == ip_1:
		cnt += 1
		if cnt > 1:
			acc -= int(input_data[ip_1][1])
			break
	ip += 1	# step the instruction pointer
	print(" ".join(input_data[ip]), acc, cnt, sep="   ")

print(f"Accumulator value for the 1st iteration: {acc}")
