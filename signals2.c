#include <stdio.h>
#include <signal.h>
#include <stdlib.h>
#include <unistd.h>

void signal_handler(int signo) {
	if (signo == SIGRTMIN) {
		printf("Recebido sinal de teste 1\n");
	}
	if (signo == SIGRTMIN+1) {
		printf("Recebido sinal de teste 2\n");
	}
	if (signo == SIGRTMIN+2) {
		printf("Recebido sinal para terminar execução\n");
		kill(getpid(), SIGTERM);
	}
}

int main(int argc, char *argv[]) {

	if (argc != 2) {
		printf("Formato: ./receive_signal wait_format(1 = busy; 2 = blocking) \n");
		return 0;
	} 
  
	signal(SIGRTMIN, signal_handler); // Sinal 34
	signal(SIGRTMIN+1, signal_handler); // Sinal 35
	signal(SIGRTMIN+2, signal_handler); // Sinal 36

	printf("Número do processo: %d\n", getpid());

	volatile int i = 0;

	while(1) {
		printf("Aguardando..\n");
		
		if (atoi(argv[1]) == 1) {
			while(i == 0) {} // Busy wait
		}

		if (atoi(argv[1]) == 2) {
			sleep(5); // Blocking wait
		}

		else {
			printf("Modo de bloqueio não disponível\n");
			return 1;
		}		
	}
	
	return 0;
}