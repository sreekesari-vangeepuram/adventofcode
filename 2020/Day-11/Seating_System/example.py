#!/usr/bin/env python

from copy import deepcopy

def get_diagonal_indices(row, col, max_row, max_col):
  tmp_list = list()
  for i in range(row-1, (row+1) + 1):
    for j in range(col-1, (col+1) + 1):
      if i >= 0 and j >= 0:
        if i <= max_row and j <= max_col:
          tmp_list.append((i, j))
  tmp_list.remove((row, col))
  return tmp_list

def manage_seats(mat):
  new_mat = deepcopy(mat)
  for row, seats in enumerate(mat):
    for col, seat in enumerate(seats):
      tmp = list(map(lambda t: mat[t[0]][t[1]], get_diagonal_indices(row, col, len(mat)-1, len(mat[0])-1))).count('#')
      if tmp >= 4:
        if mat[row][col] == "#":
          new_mat[row][col] = "L"
      elif tmp == 0:
        if mat[row][col] == "L":
          new_mat[row][col] ="#"
  return new_mat

seat_matrix = open("sample.txt").read().strip().split("\n")
initial_filling = lambda seats: [[seat for seat in row.replace("L", "#")] for row in seats]
buf_matrix =  manage_seats(initial_filling(seat_matrix))

while True:
	tmp = buf_matrix
	buf_matrix = manage_seats(tmp)
	if tmp == buf_matrix:
		print("\n".join(["".join(row) for row in buf_matrix]).count("#"))
		break
