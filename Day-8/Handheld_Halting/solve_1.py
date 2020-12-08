#!/usr/bin/env python

input_data = list(map(lambda x: x.split(), open("input.txt").read().strip().split("\n")))

acc = 0 # accumulator holding buffer variable
ip = 0 # instruction pointer
tracker = [False]*len(input_data) # `ip` recorder

print("Instructions".ljust(15),"\t|\t", "Accumulator")
print("-"*15,"\t|\t","-"*15)
while not tracker[ip]:
	tracker[ip] = True
	if input_data[ip][0] == "acc":
		acc += int(input_data[ip][1])
	if input_data[ip][0] == "jmp":
		ip += int(input_data[ip][1]) - 1 # `-1` to balance the ip increment in last step
	ip += 1	# step the instruction pointer
	print(" ".join(input_data[ip]).ljust(15),"\t|\t", acc)

print(f"Accumulator value: {acc}")
