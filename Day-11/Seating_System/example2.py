#!/usr/bin/env python

from copy import deepcopy

def get_occupied_in_sight(matrix, row, col):
  adj_pos = [
           (-1, 1),  (0, 1),  (1, 1),
           (-1, 0),           (1, 0),
           (-1, -1), (0, -1), (1, -1),
  ]
  counter = 0
  for i, j in adj_pos:
    x, y = row + i, col + j
    while 0 <= x <= len(matrix)-1 and 0 <= y <= len(matrix[0])-1:
      if matrix[x][y] == "L": break
      if matrix[x][y] == "#":
        counter += 1
        break
      x += i; y += j
  return counter

seat_matrix = list([let for let in line] for line in open("sample.txt").read().strip().split("\n"))

mat = deepcopy(seat_matrix)
while True:
  prev_mat = deepcopy(mat)
  for i, row in enumerate(prev_mat):
    for j, col in enumerate(row):
      if col == ".": continue
      occupied = get_occupied_in_sight(prev_mat, i, j)
      if col == "L" and occupied == 0:
        mat[i][j] = "#"
      elif col == "#" and occupied >= 5:
        mat[i][j] = "L"
  if mat == prev_mat:
    break
  prev_mat = mat

print(sum(row.count("#") for row in mat))
