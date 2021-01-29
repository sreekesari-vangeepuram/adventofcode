from itertools import combinations

items = [
    "space heater",
    "antenna",
    "whirled peas",
    "manifold",
    "dark matter",
    "spool of cat6",
    "bowl of rice",
    "klein bottle",
]

# Initially, drop-all!
for item in items:
    print(f"drop {item}")

for i in range(2**len(items)):

    buff_set = set()

    for n in range(len(items)):
        if (1 & (i >> n)) == 1:
            buff_set.add(items[n])
    print(i, buff_set, "\n\n")

    for item in buff_set:
        print(f"take {item}")
    print("north")

    for item in buff_set:
        print(f"drop {item}")


