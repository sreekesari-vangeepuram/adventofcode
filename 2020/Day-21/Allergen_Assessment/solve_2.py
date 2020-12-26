#!/usr/bin/env python

def pair_up(ingredients_list):
    buff_dict = dict()
    all_ingredients = list()

    for row in ingredients_list:
        ingredients, allergens = row.replace(")", "").split(" (contains ")
        ingredients, allergens = set(ingredients.split()), set(allergen.strip() for allergen in allergens.split(","))
        all_ingredients += list(ingredients)
        for allergen in set(allergens):
            buff_dict[allergen] = buff_dict.get(allergen, ingredients).intersection(ingredients)

    return buff_dict, all_ingredients


ingredients_list = open("input.txt").read().strip().split("\n")
pairs, all_ingredients = pair_up(ingredients_list)

verified_allergens, verified_ingredients = set(), set()
cdil = list() # canonical dangerous ingredients list -> cdil
while len(pairs.keys()) != 0:
    for allergen, ingredients in pairs.items():
        if len(ingredients) == 1:
            verified_allergens.add(allergen)
            ingredient = ingredients.pop()
            cdil += [(allergen, ingredient)]
            verified_ingredients.add(ingredient)
        else:
            pairs[allergen] = ingredients - verified_ingredients

    for allergen in verified_allergens:
        if allergen in pairs.keys():
            _ = pairs.pop(allergen)

print(f"Canonical dangerous ingredients list: {','.join(ingredient for _, ingredient in sorted(cdil))}")
