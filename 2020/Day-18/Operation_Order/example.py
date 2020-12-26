#!/usr/bin/env python

import re

# Operators = {+, *, ()}
# Operator precedence = (), left-to-right

def edit_expression(exp_list):
    buffer_string = ""
    for i, element in enumerate(exp_list):
        if type(element) == type(list()):
            buffer_string += "(" + edit_expression(element) + ")"
        elif element.isdigit():
            buffer_string = "(" + buffer_string  + element + ")"
        else:
            buffer_string += element
    return buffer_string

def tokenize(string):
    #_tokenizer = re.compile(r'\s*([()])\s*').split  # with    spaces
    _tokenizer = re.compile(r'(?:([()])|\s+)').split # without spaces
    return filter(None, _tokenizer(string))

def parse_as_list(expr):
    stack = []  # or a `collections.deque()` object, which is a little faster
    top = items = []
    for token in tokenize(expr):
        if token == '(':
            stack.append(items)
            items.append([])
            items = items[-1]
        elif token == ')':
            if not stack:
                raise ValueError("Unbalanced parentheses")
            items = stack.pop()
        else:
            items.append(token)
    if stack:
        raise ValueError("Unbalanced parentheses")
    return top

# Source of `tokenize` and `parse_conditions`:
# https://stackoverflow.com/questions/54959875/recursive-parentheses-parser-for-expressions-of-strings

expressions = [expression.strip() for expression in open("sample.txt")]

accumulator = 0
for expression in expressions:
    accumulator += eval(edit_expression(parse_as_list(expression)))

print(f"Sum of the resulting values: {accumulator}")
