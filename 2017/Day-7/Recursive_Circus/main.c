#include <alloca.h>
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX_LEN 10
#define MAX_CHILD 10

typedef enum {
	false,
	true
} bool;

typedef char *string;

typedef struct program {
	string name;           // Name of the program
	int    weight;         // Weight of the program
	bool   has_child;      // To decide whether
	struct program *child[MAX_CHILD]; // (Linked) List of (child) programs
} Program;


typedef struct Node {
	Program *node;
	struct Node *previous;
	struct Node *next;
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

	// Initialize a Tower
	Tower *programs, *tbuff1, *tbuff2;
	programs = (Tower *) alloca(sizeof(Tower));
	tbuff1 = (Tower *) alloca(sizeof(Tower));
	tbuff2 = (Tower *) alloca(sizeof(Tower));

	programs->node = (Program *) alloca(sizeof(Program));
	programs->previous = NULL;
	programs->next = NULL;

	// Initialize buffers
	tbuff1 = tbuff2 = programs;
	Program *pbuff;

	char c, sbuff[MAX_LEN];
	int I = 0, C = 0;

	while ((c = fgetc(fhand)) != EOF)
		switch (c) {
		case 0x20: // Space character
			sbuff[I] = 0x0; // Suffix the null character
			tbuff2->node->name = (string) alloca(sizeof(strlen(sbuff) + 1)); 
			strcpy(tbuff2->node->name, sbuff);

			c = fgetc(fhand); // Read the character after 0x20
			if (c == 0x28) // Open round bracket character
			{
				I = 0;
				// Read the weight of the current program
				// ... until 0x29 (closed round bracket character)
				while ((c = fgetc(fhand)) != 0x29)
					sbuff[I++] = c;

				tbuff2->node->weight = stoi(sbuff);
			}

			c = fgetc(fhand); // Read the character after 0x29
			switch (c) {
			case 0xA: // Newline character
				tbuff2->node->has_child = false;
				tbuff2->node->child = NULL;
				goto LINKING;

			case  0x2D: // Hyphen character
				tbuff2->node->has_child = true;
				for (C = 0; C < MAX_CHILD; ++C) // Alloc memory to children
				{
					tbuff2->node->child[C] =
						(Program *) alloca(sizeof(Program));
				}

				c = fgetc(fhand); // Skip `>`  (0x3E) character
				c = fgetc(fhand); // Skip `\s` (0x20) character

				I = 0, C = 0;
				pbuff = tbuff2->node->child;
				while ((c = fgetc(fhand)) != 0xA)
				{
					if (c == 0x2C) // Comma character
					{
						sbuff[I] = 0x0; // Null character

						(pbuff + C)->name =
							(string) alloca(sizeof(strlen(sbuff) + 1));

						strcpy((pbuff + C)->name, sbuff); // Copy the name
						(pbuff + C)->weight = -1;         // Cleaning out garbage
						(pbuff + C)->has_child = false;   // Temporarily
						(pbuff + C)->child[C++] = NULL;

						c = fgetc(fhand); // Skip space (0x20) character
						I = 0;

						continue;
					}	

					sbuff[I++] = c;
				}

				// Read the last child of the current program
				strcpy((pbuff + C)->name, sbuff); // Copy the name
				(pbuff + C)->weight = -1;         // Cleaning out garbage
				(pbuff + C)->has_child = false;   // Temporarily
				(pbuff + C)->child[C] = NULL;

				goto LINKING;
			}
		
			LINKING:	
				tbuff1->next = (Tower *) alloca(sizeof(Tower));
				tbuff1->next->previous = tbuff2;
				tbuff1->next->next = NULL;
				tbuff1 = tbuff2;

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
