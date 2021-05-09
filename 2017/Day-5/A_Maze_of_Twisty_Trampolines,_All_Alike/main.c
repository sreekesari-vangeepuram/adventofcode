#include <stdio.h>
#include <stdlib.h>

void run_program(int *, int *);
void duplicate_ints(int *, int *, int);
int line_count(FILE *);
int stoi(const char *);

int N = 0, part[2];

int main(int argc, char *argv[])
{

	if (argc < 2)
	{
		printf("\e[31;1;4mUsage\e[0m: \e[32m./main\e[0m \e[34m<input-file>\e[0m\n");
		return 1; // NO INPUT FILE
	}

	FILE *fhand = fopen(argv[1], "r");

	N = line_count(fhand);
	int instructions[N];

	char c, sbuff[10];
	int I = 0, J = 0;
	while((c = fgetc(fhand)) != EOF)
		switch (c) {
		case 0xA: // Newline character
			sbuff[I] = 0x0;
			instructions[J++] = stoi(sbuff);
			I = 0; break;

		default:
			sbuff[I++] = c;
		}

	run_program(instructions, part);

	printf("\e[33mNumber of steps to reach the \e[31mexit\e[0m [CASE \e[36mA\e[0m]: \e[2m%d\e[0m\n", part[0]);
	printf("\e[33mNumber of steps to reach the \e[31mexit\e[0m [CASE \e[36mB\e[0m]: \e[2m%d\e[0m\n", part[1]);
	
	fclose(fhand);
	
	return 0; // SUCCESS
}

void run_program(int *instructions, int *part)
{
	int *ins = malloc(N * sizeof(int));
	duplicate_ints(instructions, ins, N);

	int ip = 0, jmp = 0, steps = 0;
	while (ip < N)
	{
		jmp = ins[ip]++;
		ip += jmp;
		steps++;
	}

	part[0] = steps;

	steps = jmp = ip = 0;
	while (ip < N)
	{
		jmp = instructions[ip];

		if (jmp >= 0x3)
			instructions[ip]--;
		else
			instructions[ip]++;

		ip += jmp;
		steps++;
	}

	part[1] = steps;

	return;
}

void duplicate_ints(int *src, int *dest, int N)
{
	for (int i = 0; i < N; ++i)
		dest[i] = src[i];

	return;
}

int line_count(FILE *fhand)
{
	rewind(fhand);

	char c;
	int lines = 0;
	while ((c = fgetc(fhand)) != EOF)
		if (c == 0xA) // Newline character
			lines++;

	rewind(fhand);
	return lines;
}

int stoi(const char *str)
{
	int ibuff = 0, i = 0, _ve = 1;
	while (str[i] != 0x0)
	{
		if (str[i] == 0x2D) // Negative sign
		{
			_ve = -1;
			++i; continue;
		}

		if (!ibuff)
			ibuff = ((int)str[i]) - 0x30;
		else
		{
			ibuff *= 10;
			ibuff += ((int)str[i]) - 0x30;
		}

		++i;
	}

	return _ve * ibuff;
}
