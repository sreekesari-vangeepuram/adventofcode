#include <stdio.h>
#define MAX_STR 10

int ROWS = 0, COLS = 0;
void generate_checksum(int *, FILE *);

int stoi(const char *str)
{
	int ibuff = 0, i = 0;
	while (str[i])
	{
		if (!ibuff)
			ibuff = ((int)str[i]) - 0x30;
		else
		{
			ibuff *= 10;
			ibuff += ((int)str[i]) - 0x30;
		}

		i++;
	}

	return ibuff;
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
	int ibuff = 0;

	// Get max. row and column count
	while ((c = fgetc(fhand)) != EOF)
		switch (c) {
		case 0x9: // Tab character
			ibuff++;
			break;

		case 0xA: // Newline character
			if (ibuff > COLS)
				COLS = ++ibuff;

			ibuff = 0;
			ROWS++;
			break;
		}

	rewind(fhand);

	int part[2];
	generate_checksum(part, fhand);
	
	printf("\033[32mOLD CHECKSUM\033[0m: \033[34m%d\033[0m\n", part[0]);
	printf("\033[32mNEW CHECKSUM\033[0m: \033[34m%d\033[0m\n", part[1]);
	
	return 0; // SUCCESS
}

void generate_checksum(int *part, FILE *fhand)
{
	int spreadsheet[ROWS][COLS];

	char c, str[MAX_STR];
	int i = 0, j = 0, ip = 0;

	// Parse data into spreadsheet [2x2]
	while ((c = fgetc(fhand)) != EOF)
		switch (c) {
		case 0x9: // Tab character
			str[ip] = 0x0; // Null character
			spreadsheet[i][j++] = stoi(str);
			ip = 0; break;

		case 0xA: // Newline character
			str[ip] = 0x0; // Null character
			spreadsheet[i++][j] = stoi(str);
			j = ip = 0; break;

		default: str[ip++] = c;
		}

	part[0] = 0;
	part[1] = 0;

	int min, max, k, rem;
	for (i = 0; i < ROWS; ++i)
	{
		min = spreadsheet[i][0];
		max = spreadsheet[i][0]; 

		for (j = 0; j < COLS; ++j)
		{
			if (min > spreadsheet[i][j])
				min = spreadsheet[i][j];
			if (max < spreadsheet[i][j])
				max = spreadsheet[i][j];

			while (k < COLS)
			{
				if (!(spreadsheet[i][k] % spreadsheet[i][j]) && (j != k))
					if (spreadsheet[i][k] > spreadsheet[i][j])
					rem = spreadsheet[i][k] / spreadsheet[i][j];
				++k;
			}

			k = 0;
		}

		part[0] += max - min;
		part[1] += rem;
	}
 	
	fclose(fhand);

	return;
}
