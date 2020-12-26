#!/usr/bin/env python

# card_pubkey, door_pubkey = map(int, open("input.txt").read().strip().split("\n"))
# card_pubkey, door_pubkey = 13316116, 13651422 # Obtained from `sample.txt`
# prime = 20201227

# ****************************************************** #
# Hard-coding numbers insted of assigning to variables   #
# decreases the runtime due to decrease in IO-operations #
# ****************************************************** #

subnum, enckey = 7,  13651422 # door_pubkey
while subnum!= 13316116: # card_pubkey
    enckey = (enckey * 13651422) % 20201227     # Encryption key
    subnum = (subnum * 7) % 20201227            # Subject number

print(f"Encryption key: {enckey}")
