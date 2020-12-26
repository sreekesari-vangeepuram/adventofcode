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


ingredients_list = open("sample.txt").read().strip().split("\n")
pairs, all_ingredients = pair_up(ingredients_list)

verified_allergens, verified_ingredients = set(), set()
while len(pairs.keys()) != 0:
    for allergen, ingredients in pairs.items():
        if len(ingredients) == 1:
            verified_allergens.add(allergen)
            verified_ingredients.add(ingredients.pop())
        else:
            pairs[allergen] = ingredients - verified_ingredients

    for allergen in verified_allergens:
        if allergen in pairs.keys():
            _ = pairs.pop(allergen)

unmatched_ingredients = set(all_ingredients) - verified_ingredients
appearances = sum(all_ingredients.count(ingredient) for ingredient in unmatched_ingredients)

print(f"Count of the [duplicate] unmatched ingredinets: {appearances}")
