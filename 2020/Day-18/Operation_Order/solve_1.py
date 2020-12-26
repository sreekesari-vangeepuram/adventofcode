#!/usr/bin/env python

def eval_expr(expr: str) -> int:
    constants, operators = list(), list()

    def smash_stack() -> None:
        t1, t2 = constants.pop(), constants.pop()
        operator = operators.pop()
        constants.append(eval(f"{t1}{operator}{t2}"))

    def precedence_fn(operator: str) -> int:
        if operator in {"+", "*"}:
            return 1
        return 0

    for char in "".join(expr.split()):
        if char.isdigit():
            constants.append(int(char))
        elif char == "(":
            operators.append(char)
        elif char == ")":
            while len(operators) > 0 and operators[-1] != "(":
                smash_stack()
            operators.pop()
        elif char in {"+", "*"}:
            while len(operators) > 0 and precedence_fn(operators[-1]) >= precedence_fn(char):
                smash_stack()
            operators.append(char)

    while len(operators):
        smash_stack()

    return constants[-1]

# Reference(s) :
# http://www.openbookproject.net/books/pythonds/BasicDS/InfixPrefixandPostfixExpressions.html

values = [eval_expr(expr.strip()) for expr in open("input.txt")]

print(f"Sum of resulting values: {sum(values)}")
