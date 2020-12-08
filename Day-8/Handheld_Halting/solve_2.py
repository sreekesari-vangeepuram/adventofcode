#!/usr/bin/env python

input_data = list(map(lambda x: x.split(), open("input.txt").read().strip().split("\n")))

def get_accumulator(ins=input_data):
	acc = 0 # accumulator holding buffer variable
	ip = 0 # instruction pointer
	tracker = [False]*len(ins) # `ip` recorder
	print("_"*15,"\t|\t","_"*15)
	print("Instructions".ljust(15),"\t|\t", "Accumulator")
	print("-"*15,"\t|\t","-"*15)
	while ip in range(len(ins)) and not tracker[ip]:
		tracker[ip] = True
		if ins[ip][0] == "acc": acc += int(ins[ip][1])
		if ins[ip][0] == "jmp": ip += int(ins[ip][1])-1
		ip += 1
		try:
			print(" ".join(ins[ip]).ljust(15),"\t|\t", acc)
		except:
			print("_"*40)
	return acc if ip == len(ins) else "Stuck in <Infinite-Loop>"

op_shift = {"nop": "jmp", "acc": "acc", "jmp": "nop"}

for i in range(len(input_data)):
	input_data[i][0] = op_shift[input_data[i][0]]
	acc = get_accumulator(input_data)
	input_data[i][0] = op_shift[input_data[i][0]]
	if acc != "Stuck in <Infinite-Loop>":
		break

print(f"Accumulator value: {acc}")
