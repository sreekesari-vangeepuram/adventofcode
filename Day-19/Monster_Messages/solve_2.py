#!/usr/bin/env python

def check_match(message, rule0):
    matched = True
    if len(rule0) == 0 or len(message) == 0:
        return len(rule0) == 0 and len(message) == 0
    elif len(rule0) > len(message):
        return not matched

    tmp = rule0.pop()
    if isinstance(tmp, str):
        if message[0] == tmp:
            return check_match(message[1:], rule0.copy())
    else:
        for rule in customised_rules[tmp]:
            if check_match(message, rule0 + list(reversed(rule))):
                return matched
    return not matched


def customise_rules(rules):
    buffer_dict = dict()
    for rule in rules:
        key, values = rule.split(": ")
        if values[0] == '"':
            buffer_dict[int(key)] = values.replace('"', '')
        else:
            buffer_list = list()
            values = values.split(" | ")
            for value in values:
                buffer_list.append(list(int(v) for v in value.split()))
            buffer_dict[int(key)] = buffer_list

    buffer_dict[8] = [[42], [42, 8]]
    buffer_dict[11] = [[42, 31], [42, 11, 31]]
    return buffer_dict


given_rules, given_messages = open("input.txt").read().split("\n\n")
given_rules, given_messages = given_rules.strip().split("\n"), given_messages.strip().split("\n")
customised_rules = customise_rules(given_rules)

accumulator = 0
for message in given_messages:
    # Graph + Stack : `customised_rules`
    if check_match(message, customised_rules[0][0][::-1]):
        accumulator += 1

print(f"Messages mathing `Rule-0`: {accumulator}")
