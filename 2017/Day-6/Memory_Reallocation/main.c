#include <stdio.h>
#include <stdlib.h>

#define N 16

typedef enum {
	false,
	true
} bool;

typedef struct Bank memory_t;
struct Bank {
	int cycle_id;
	int sector[N];
	memory_t *next;
};

//// Side effects <=> Special effects ////
void realloc_memory(int *);
void update_memory(int *);
bool match(memory_t *, memory_t *);
int stoi(const char *);

int part[2] = {0};

int main(int argc, char *argv[])
{

	if (argc < 2)
	{
		printf("\e[31;1;4mUsage\e[0m: \e[32m./main\e[0m \e[34m<input-file>\e[0m\n");
		return 1; // NO INPUT FILE
	}

	FILE *fhand = fopen(argv[1], "r");

	/// PRE INTIAL STATE ///
	int mem[N] = {0};

	char c, sbuff[10];
	int I = 0, J = 0;
	while ((c = fgetc(fhand)) != EOF) // Until `c` hit a newline character
		switch (c) {
		case 0x9: // Tab character
			sbuff[I] = 0x0; // Null character
			if (J < N) mem[J++] = stoi(sbuff);
			I = 0; break;

		case 0xA: // Newline character [For last number stored in `sbuff`]
			sbuff[I] = 0x0; // Null character
			if (J < N) mem[J++] = stoi(sbuff);
			break;

		default:
			// Add to buffer
			sbuff[I++] = c;
		}

	realloc_memory(mem);

	printf("\e[33mRedistribution cycle count \e[0m: \e[2m%d\e[0m cycles\n", part[0]);
	printf("\e[33mNext same configuration at \e[0m: \e[2m%d\e[0m cycles\n", part[1]);

	fclose(fhand);
	
	return 0; // SUCCESS
}

void realloc_memory(int *mem)
{

	/// Initial storage space state ///
	memory_t *storage_space = (memory_t *) malloc(sizeof(memory_t));
	storage_space->cycle_id = 0;
	for (int b = 0; b < N; ++b)
		storage_space->sector[b] = mem[b];

	// `mem` itself used as buffer
	// whereas `buffer_space` of type <memory_t>

	// `buffer_space` is used to copy the
	// state of `mem` into `storage_space`
	memory_t *buffer_space = (memory_t *) malloc(sizeof(memory_t));

	// `mapping_space` is used to connect
	// next `buffer_space` with old `buffer_space`
	memory_t *mapping_space = NULL;

	storage_space->next = buffer_space;

	bool found = false;
	int cycle_id = 0, i;

	///// Update state /////
	while (!found)
	{
		buffer_space->cycle_id = ++cycle_id;

		////////////////// UPDATE MEMORY //////////////////

		update_memory(mem);

		for (i = 0; i < N; ++i) // Update buffer sectors
			buffer_space->sector[i] = mem[i];

		found = match(storage_space, buffer_space);
	
		////////// LOG BUFFER INTO STORAGE SPACE //////////

		buffer_space->next = NULL;
		mapping_space = buffer_space;

		// Shift buffer space to new location
		// during each cycle or memory reallocation

		buffer_space = (memory_t *) malloc(sizeof(memory_t));
		mapping_space->next = buffer_space;
	}	

	// Free used memory
	free(storage_space);
	free(buffer_space);
	free(mapping_space);

	return;
}

void update_memory(int *mem)
{
	// Find the bank with
	// maximum blocks!
	int max_b = 0, I = 0;
	for (int i = 0; i < N; ++i)
	{	
		// Tie always won by
		// bank near to Bank - 1!
		if (max_b < mem[i])
		{
			max_b = mem[i];
			I = i;
		}
	}

	// Remove all the blocks from
	// the bank with maximum blocks
	mem[I++] = 0;

	// Distribute all the blocks
	// collected from the big bank
	while (max_b)
	{
		if (I == N) I = 0;
		mem[I++]++; // Start adding a block from the next bank
		--max_b;
	}

	return;
}

bool match(memory_t *storage_space, memory_t *buffer_space)
{
	int i, acc;
	bool found = false;

	// Check all the states except
	// the last one <=> buffer_space

	while (storage_space->next != NULL)
	{
		acc = 0;
		for (i = 0; i < N; ++i)
			if (buffer_space->sector[i] == storage_space->sector[i])
				++acc;

		if (acc == N)
		{
			found = true;
			part[0] = buffer_space->cycle_id;

			// Next similar state with same frequency!
			part[1] = buffer_space->cycle_id - storage_space->cycle_id;
		}

		// Move to the next state
		storage_space = storage_space->next;
	}

	return found;
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
