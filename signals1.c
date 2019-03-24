#include <stdio.h>
#include <signal.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
  
	if (argc != 3) {
		printf("Formato: send_signal pid signal\n");
		return 0;
	} 

	int process_exists = kill(atoi(argv[1]), 0);

	if (process_exists != 0) {
		printf("Processo selecionado não existe\n");
		return 1;
	}

	else {
		int success = kill(atoi(argv[1]), atoi(argv[2]));
		if(success == 0){
			printf("Sinal enviado com sucesso\n");
		}
	}

	return 0;
}