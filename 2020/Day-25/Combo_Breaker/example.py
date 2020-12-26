#!/usr/bin/env python

# card_pubkey, door_pubkey = map(int, open("input.txt").read().strip().split("\n"))
# card_pubkey, door_pubkey = 5764801, 17807724 # Obtained from `sample.txt`
# prime = 20201227

# ****************************************************** #
# Hard-coding numbers insted of assigning to variables   #
# decreases the runtime due to decrease in IO-operations #
# ****************************************************** #

subnum, enckey = 7, 17807724
while subnum!= 5764801:
    enckey = (enckey * 17807724) % 20201227     # Encryption key
    subnum = (subnum * 7) % 20201227            # Subject number

print(f"Encryption key: {enckey}")
