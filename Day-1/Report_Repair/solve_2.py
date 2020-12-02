#!/usr/bin/env python

pairs = list()
with open('input.txt') as input_data:
        data = list(map(int, input_data))
        for i in data:
                for j in data:
                        for k in data:
                                if i + j + k == 2020:
                                        print(i*j*k)
