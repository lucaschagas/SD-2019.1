#include <stdio.h> 
#include <stdlib.h> 
#include <unistd.h> 
#include <sys/types.h> 
#include <string.h> 
#include <sys/wait.h> 
#include <time.h>

#define SIZE 21

int random_value(void) {
	srand(time(NULL));
	return (rand() % 99) + 1;
}

void is_prime(int n) {

    int i, flag = 0;

    for(i = 2; i <= n/2; ++i) {
        if (n%i == 0) {
        	flag = 1;
        	break;
        }
    }

    if (n == 1) {
    	printf("1 não é primo.\n");
    }

    else {
        if (flag == 0) {
        	printf("%d é primo.\n", n);
        }
        else {
        	printf("%d não é primo.\n", n);
        }
    }
}

int main(int argc, char *argv[]) {

	if (argc != 2) {
		printf("Formato: ./pipes quantidade_de_valores \n");
		return 0;
	}

	int pipe1[2];

	int p = pipe(pipe1);
	int f = fork();

    // Processo pai
    if (f > 0) {
    	char valor_enviado[SIZE];
    	int valor = 1;
    	int contador = 0;
    	close(pipe1[0]);

    	while (contador < atoi(argv[1])) {
    		int valor_aleatorio = random_value();
    		valor += valor_aleatorio;
    		sprintf(valor_enviado, "%d", valor);
    		write(pipe1[1], valor_enviado, SIZE);

    		contador++;
    	}

    	write(pipe1[1], "0", SIZE);
    	close(pipe1[1]);
    	exit(0);
    } 

    // Processo filho
    else {
        char valor_recebido[SIZE];
    	int valor_final;
    	close(pipe1[1]);

    	while (valor_final != 0) {
    		read(pipe1[0], valor_recebido, SIZE);
    		valor_final = atoi(valor_recebido);

    		if (valor_final > 0) {
    			is_prime(valor_final);
    		}
    	}

        close(pipe1[0]);
    	exit(0);
    }

	return 0;
}