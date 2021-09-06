#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX_LEN 10

typedef enum {
	false,
	true
} bool;

typedef char *string;

typedef struct program {
	string name;           // Name of the program
	int    weight;         // Weight of the program
	bool   has_child;      // To decide whether
	struct program *child; // (Linked) List of (child) programs
} Program;

typedef struct Node {
	Program node;
	Node *previous;
	Node *next;
} Tower;

////// FUNCTION PROTOTYPES //////

int stoi(const char *);


/////// GLOBAL VARIABLES ///////

int part[2] = {0};


int main(int argc, char *argv[])
{

	if (argc < 2)
	{
		printf("\e[31;1;4mUsage\e[0m: \e[32m./main\e[0m \e[34m<input-file>\e[0m\n");
		return 1; // NO INPUT FILE
	}

	FILE *fhand = fopen(argv[1], "r");
	assert(fhand != NULL);

	// Initialize TOWER
	Tower *programs, *tbuff;
	programs->node = (Program) malloc(sizeof(Program));
	programs->previous = NULL;
	programs->next = NULL;

	tbuff = programs; // Tower buffer

	char c, sbuff[MAX_LEN];
	int I = 0;

	while ((c = fgetc(fhand)) != EOF)
		switch (c) {
		case 0x20: // Space character
			sbuff[I] = 0x0; // Suffix the null character
			tbuff->node.name = (string) malloc(sizeof(strlen(sbuff) + 1)); 
			strcpy(tbuff->node.name, sbuff);

			c = fgetc(fhand); // Read the character after 0x20
			if (c == 0x28) // Open round bracket character
			{
				I = 0;
				while ((c = fgetc(fhand)) != 0x29) // Closed round bracket character
					sbuff[I++] = c;

				tbuff->node.weight = stoi(sbuff);

			}

			c = fgetc(fhand); // Read the character after 0x29
			switch (c) {
			case 0xA: // Newline character
				tbuff->node.has_child = false;
				tbuff->node.child = NULL;
				break;

			case  0x2D: // Hyphen character
				// Add child(ren) to the present program node

				// Also deal with comma (0x2C)
				// until newline character

				c = fgetc(fhand); // Skip `>`  (0x3E) character
				c = fgetc(fhand); // Skip `\s` (0x20) character

				I = 0;

				/* SET
				 * program.has_child = true;
				 * skip `,\s` delimiters, until `\n`
				 */
			}
			
			I = 0; break;

		default:
			sbuff[I++] = c;
			break;
		}

	//printf("\e[33m  \e[0m: \e[2m%d\e[0m\n", part[0]);
	//printf("\e[33m  \e[0m: \e[2m%d\e[0m\n", part[1]);

	fclose(fhand);
	
	return 0; // SUCCESS
}

///////////////////////////////////////////////////
// Implement a TRIE finally with the parsed data //
///////////////////////////////////////////////////


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
