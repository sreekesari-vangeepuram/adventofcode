#include <stdio.h>
#define APPROX_SIZE 3000

// Full of side-effects
void validate_captcha(int *list, int *part, const int N)
{
	const int OFF1 = N - 1, OFF2 = N / 2;
	part[0] = 0;
	part[1] = 0;

	for (int i = 0; i < N; ++i)
	{
		if (i < OFF1)
			if (list[i] == list[i + 1])
				part[0] += list[i];

		if (i + OFF2 < N)
			if (list[i] == list[i + OFF2])
				part[1] += list[i] * 2;
	}

	if (list[0] == list[OFF1])
		part[0] += list[0];

	return;
}

int main(int argc, char *argv[])
{

	if (argc < 2)
	{
		printf("\033[31;1;4mUsage\033[0m: \033[32m./main\033[0m \033[34m<input-file>\033[0m\n");
		return 1; // NO INPUT FILE
	}

	FILE *fhand = fopen(argv[1], "r");
	char c;
	int digit[APPROX_SIZE], N = 0;
	
	while((c = fgetc(fhand)) != 0xA) // Newline character
	{
		digit[N] = c - 0x30; // ASCII code for character "1" is 49
		N++; // Character count
	}

	int part[2];
	validate_captcha(digit, part, N);

	printf("\033[32;1;4mOLD CAPTCHA CODE\033[0m: \033[34m%d\033[0m\n", part[0]);
	printf("\033[32;1;4mNEW CAPTCHA CODE\033[0m: \033[34m%d\033[0m\n", part[1]);
	
	fclose(fhand);
	
	return 0; // SUCCESS
}
