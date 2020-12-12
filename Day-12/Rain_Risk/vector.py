class vector:

	def __init__(self, x, y, pointing_direction):
		self.x = x
		self.y = y
		self.direction = pointing_direction

	def get_pos(self):
		return (self.x, self.y, self.direction)

	def change_position(self, ins):
		d = self.direction
		if d == 'E':
			if ins[0] == 'E' or ins[0] == 'F': self.x += ins[1]
			elif ins[0] == 'W': self.x -= ins[1]
			elif ins[0] == 'N': self.y += ins[1]
			elif ins[0] == 'S': self.y -= ins[1]
		elif d == 'W':
			if ins[0] == 'W' or ins[0] == 'F': self.x -= ins[1]
			elif ins[0] == 'E': self.x += ins[1]
			elif ins[0] == 'N': self.y += ins[1]
			elif ins[0] == 'S': self.y -= ins[1]
		elif d == 'N':
			if ins[0] == 'N' or ins[0] == 'F': self.y += ins[1]
			elif ins[0] == 'W': self.x -= ins[1]
			elif ins[0] == 'E': self.x += ins[1]
			elif ins[0] == 'S': self.y -= ins[1]
		elif d == 'S':
			if ins[0] == 'S' or ins[0] == 'F': self.y -= ins[1]
			elif ins[0] == 'W': self.x -= ins[1]
			elif ins[0] == 'N': self.y += ins[1]
			elif ins[0] == 'E': self.x += ins[1]

	def change_direction(self, ins):
		d = self.direction
		if ins[0] == 'R':
			if d == 'E':
				if ins[1] == 90:
					self.direction = 'S'
				elif ins[1] == 180:
					self.direction = 'W'
				elif ins[1] == 270:
					self.direction = 'N'
			elif d == 'W':
				if ins[1] == 90:
					self.direction = 'N'
				elif ins[1] == 180:
					self.direction = 'E'
				elif ins[1] == 270:
					self.direction = 'S'

			elif d == 'N':
				if ins[1] == 90:
					self.direction = 'E'
				elif ins[1] == 180:
					self.direction = 'S'
				elif ins[1] == 270:
					self.direction = 'W'

			elif d == 'S':
				if ins[1] == 90:
					self.direction = 'W'
				elif ins[1] == 180:
					self.direction = 'N'
				elif ins[1] == 270:
					self.direction = 'E'

		elif ins[0] == 'L':
			if d == 'E':
				if ins[1] == 90:
					self.direction = 'N'
				elif ins[1] == 180:
					self.direction = 'W'
				elif ins[1] == 270:
					self.direction = 'S'
			elif d == 'W':
				if ins[1] == 90:
					self.direction = 'S'
				elif ins[1] == 180:
					self.direction = 'E'
				elif ins[1] == 270:
					self.direction = 'N'

			elif d == 'N':
				if ins[1] == 90:
					self.direction = 'W'
				elif ins[1] == 180:
					self.direction = 'S'
				elif ins[1] == 270:
					self.direction = 'E'

			elif d == 'S':
				if ins[1] == 90:
					self.direction = 'E'
				elif ins[1] == 180:
					self.direction = 'N'
				elif ins[1] == 270:
					self.direction = 'W'
	def manhattan_distance(self):
		return abs(self.x)+abs(self.y)
