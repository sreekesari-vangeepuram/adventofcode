#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX_LEN 15

typedef enum {
	false,
	true
} bool;

typedef char *string;

typedef struct {
	string words[MAX_LEN];
	int words_count;
	bool duplicates_policy;
	bool anagrams_policy;
} Passphrase;

int line_count(FILE *fhand)
{
	rewind(fhand);

	char c;
	int lines = 0;
	while ((c = fgetc(fhand)) != EOF)
		if (c == 0xA)
			lines++;

	rewind(fhand);
	return lines;
}

bool is_anagram(string, string);
void validate_passphrases(Passphrase *);
void count_valid_passphrases(int *, Passphrase *);

int LIST_SIZE, part[2];

int main(int argc, char *argv[])
{

	if (argc < 2)
	{
		printf("\e[31;1;4mUsage\e[0m: \e[32m./main\e[0m \e[34m<input-file>\e[0m\n");
		return 1; // NO INPUT FILE
	}

	FILE *fhand = fopen(argv[1], "r");
	LIST_SIZE = line_count(fhand);

	Passphrase list[LIST_SIZE];

	char c, sbuff[MAX_LEN];
	int I = 0, J = 0, line = 0;
	while ((c = fgetc(fhand)) != EOF)
		switch (c) {
		case 0x20: // Space character
			sbuff[I] = 0x0;
			list[line].words[J] = (string) malloc(strlen(sbuff) + 1);
			strcpy(list[line].words[J++], sbuff);
			I = 0; break;

		case 0xA: // Newline character
			sbuff[I] = 0x0;
			list[line].words[J] = (string) malloc(strlen(sbuff) + 1);
			strcpy(list[line].words[J++], sbuff);

			list[line].duplicates_policy = true; // Initially all are valid
			list[line].anagrams_policy = true; // Initially all are valid
			list[line].words_count = J;
			I = J = 0; ++line; break;

		default:
			sbuff[I++] = c;
		}

	validate_passphrases(list);
	count_valid_passphrases(part, list);

	printf("Valid passphrase count (Under DUPLICATE POLICY) : %d\n", part[0]);
	printf("Valid passphrase count (Under  ANAGRAM  POLICY) : %d\n", part[1]);

	fclose(fhand);
	
	return 0; // SUCCESS
}

bool is_anagram(string s1, string s2)
{
	// This function is valid only for
	// lower case English character-set
	// in ASCII code.

	int fs1[26] = {0}, fs2[26] = {0}, i = 0;
	
	// Frequency of all characters from s1
	while (s1[i] != 0x0) // Until null character
		fs1[s1[i++] - 0x61]++; // Alphabet line from `a`
	
	i = 0;
	
	while (s2[i] != 0x0) // Until null character
		fs2[s2[i++] - 0x61]++; // Alphabet line from `a`
	
	// Frequency check of all characters between s1, s2
	for (i = 0; i < 26; i++)
		if (fs1[i] != fs2[i])
			return false;
	
	return true;
}

void validate_passphrases(Passphrase *list)
{
	int i, j;
	for (int phrase = 0; phrase < LIST_SIZE; ++phrase)
	{
		for (i = 0; i < list[phrase].words_count; ++i)
		for (j = 0; j < list[phrase].words_count; ++j)
		{
			// Skip if both buffers has same word
			if (i == j) continue;

			// `strcmp` returns 0 if s1, s2 are same
			if (!strcmp(list[phrase].words[i], list[phrase].words[j]))
			{
				list[phrase].duplicates_policy = false;
				list[phrase].anagrams_policy = false;
				goto OUTER_LOOP; /* SKIP_SCOPE to save time */
			}

			// `is_anagram` returns 1 if s1, s2 are anagrams
			if (is_anagram(list[phrase].words[i], list[phrase].words[j]))
				list[phrase].anagrams_policy = false;
		}

		OUTER_LOOP: /* Do nothing */ ;
	}

	return;
}

void count_valid_passphrases(int *part, Passphrase *list)
{
	part[0] = part[1] = 0;
	for (int phrase = 0; phrase < LIST_SIZE; ++phrase)
	{
		part[0] += list[phrase].duplicates_policy; // If false (+ 0) is negligible
		part[1] += list[phrase].anagrams_policy; // If false (+ 0), is negligible
	}

	return;	
}
