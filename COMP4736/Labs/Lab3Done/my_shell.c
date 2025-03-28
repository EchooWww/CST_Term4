#include <stdio.h>
#include <sys/types.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/wait.h>
#include <signal.h>

#define MAX_INPUT_SIZE 1024
#define MAX_TOKEN_SIZE 64
#define MAX_NUM_TOKENS 64
#define MAX_BG_PROCESSES 64

pid_t bg_processes[MAX_BG_PROCESSES];
int bg_processes_count = 0;

void reap_bg_processes() {
	int status;
	pid_t pid;
	while((pid = waitpid(-1, &status, WNOHANG)) > 0) {
		printf("Background process %d completed with status %d\n", pid, status);
		for (int i = 0; i < bg_processes_count; i++) {
			if (bg_processes[i] == pid) {
				bg_processes[i] = bg_processes[--bg_processes_count];
				break;
			}
		}
	}
}

/**
 * Splits the string by space and returns the array of tokens.
 */
char **tokenize(char *line) {
	char **tokens = (char **)malloc(MAX_NUM_TOKENS * sizeof(char *));
	char *token = (char *)malloc(MAX_TOKEN_SIZE * sizeof(char));
	int i, tokenIndex = 0, tokenNo = 0;

	for (i = 0; i < strlen(line); i++) {
		char readChar = line[i];

		if (readChar == ' ' || readChar == '\n' || readChar == '\t') {
			token[tokenIndex] = '\0';
			if (tokenIndex != 0) {
				tokens[tokenNo] = (char*)malloc(MAX_TOKEN_SIZE * sizeof(char));
				strcpy(tokens[tokenNo++], token);
				tokenIndex = 0; 
			}
		}
		else {
			token[tokenIndex++] = readChar;
		}
	}

	free(token);
	tokens[tokenNo] = NULL ;
	return tokens;
}

int main(int argc, char* argv[]) {
	char line[MAX_INPUT_SIZE];            
	char **tokens;              
	int i;

	while (1) {
		/* BEGIN: TAKING INPUT */
		bzero(line, sizeof(line));
		printf("$ ");
		scanf("%[^\n]", line);
		getchar();

		/* END: TAKING INPUT */

		reap_bg_processes();
		line[strlen(line)] = '\n'; //terminate with new line
		tokens = tokenize(line);
		int num_tokens = sizeof(tokens) / sizeof(tokens[0]);

		if (tokens[0] == NULL) {
			free(tokens);
			continue;
		}
		if (strcmp(tokens[0], "exit") == 0) {
			for (i = 0; tokens[i] != NULL; i++) {
				kill(bg_processes[i], SIGTERM);
				printf("Killed process %d\n", bg_processes[i]);
			}
			free(tokens);
			printf("Exiting shell...\n");
			exit(0);
		}

		if (strcmp(tokens[0], "cd") == 0) {
			if (tokens[1] == NULL || chdir(tokens[1]) < 0) {
				printf("Shell: Incorrect directory\n");
			}
			continue;
		}

		int bg = 0;
		int i;
		for (i = 0; tokens[i] != NULL; i++);
		if (i > 0 && strcmp(tokens[i - 1], "&") == 0) {
				bg = 1;
				free(tokens[i - 1]);
				tokens[i - 1] = NULL;
		}

		pid_t pid = fork();
		if (pid == -1) {
			printf("Shell: Failed to fork\n");
			continue;
		}
		else if (pid == 0) {
			// if execvp returns -1, then the command is incorrect
			if (execvp(tokens[0], tokens) < 0) {
				printf("Shell: Incorrect command\n");
				exit(1);
			}
		}
		else {
			if (bg) {
				if (bg_processes_count<MAX_BG_PROCESSES) {
					bg_processes[bg_processes_count++] = pid;
				} 
			} else {
				int status;
				while (waitpid(pid, &status, WNOHANG) == 0) {
						reap_bg_processes();  
						usleep(100000);
				}				
			}
		}
   
		// Freeing the allocated memory	
		for(i = 0; tokens[i] != NULL; i++) {
			free(tokens[i]);
		}
		free(tokens);

	}
	return 0;
}
