#include <stdio.h>

// Input  vvvvvv
#define N 277678
#define MAX_LEN 100

int delta[8][2] = {
	{-1,  1}, {0,  1}, {1,  1},
	{-1,  0},          {1,  0},
	{-1, -1}, {0, -1}, {1, -1},
};

typedef struct {
	int value;
	int x, y;
} sector_t;

int part[2]; // Solution buffer

int abs(int);
void set_memory(int);

int main(void)
{
	set_memory(0xA); // MEMORY ALLOCATION MODE = 0xA
	
	printf("\e[32mManhattan Distance\e[0m from ");
	printf("\e[2msector\e[0m[1] to \e[2msector\e[0m[%d]: \e[1m%d\e[0m\n",
			N, part[0]);
	
	set_memory(0xB); // MEMORY ALLOCATION MODE = 0xB
	
	printf("\e[33mFirst value written\e[0m i.e. \e[31m> %d\e[0m: \e[1m%d\e[0m\n",
			N, part[1]);
	
	return 0; // SUCCESS
}

void set_memory(int mode)
{
	int size = (mode == 0xA) ? N : MAX_LEN; // CHOOSE MEMORY SIZE W.R.T. `mode`
	sector_t mem[size]; // ALLOCATING MEMORY BUFFER OF SIZE `size`

	// Clean garbage data
	for (int I = 0; I < size; ++I)
		mem[I].value = 0;

	int x = 0, y = 0,
		dx = 1, dy = 0,
		z = 1, sector = 0,
		ibuff = 0,
		k = 0, J = 0;

	for (int i = 0; i < N; ++i)
	{
		mem[i].x = x;
		mem[i].y = y;

		switch (mode) {
		case 0xA:
			mem[i].value = i + 1;
			if (i == N - 1)
			{
				part[0] = abs(mem[N - 1].x) + abs(mem[N - 1].y);
				return;
			}
			break;

		case 0xB:
			// Starting case 
			if ((x == 0) && (y == 0))
			{
				mem[i].value = 1;
				break;
			}

			ibuff = 0;
			for (J = 0; J < MAX_LEN; ++J)
			for (k = 0; k < 8; ++k)
			if ((x + delta[k][0] == mem[J].x))
			if ((y + delta[k][1] == mem[J].y))
				ibuff += mem[J].value;
		
			mem[i].value = ibuff;
			break;

		default:
			printf("An ERROR encountered!\n");
			return;
		}
	
		if (mem[i].value > N)
		{
			part[1] = mem[i].value;
			return;	
		}

		x += dx; y += dy;

		if (++sector == z)
		{
			sector = 0; // Passed!

			// Rotate sprially
			// at corner cases

			ibuff = dx;
			dx = -dy;
			dy = ibuff;

			if (!dy) z++; // Move to next layer...
		}
	}

	return;
}

int abs(int num)
{
	if (num < 0) return -num;
	else return num;
}

