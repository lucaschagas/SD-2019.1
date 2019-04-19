#include <pthread.h>
#include <stdio.h>
#include <stdlib.h> 
#include <time.h> 
#include <stdint.h>

// gcc spinlock.c -pthread -o spinlock

int lock = 0;
int acumulador = 0;
int indice_temp = 0;
int parcelas = 0;
int resto = 0;

void acquire();

void release();

int random_value(){
	int valor = (rand() % 201) - 100;
	return valor;
}

void *func1(int8_t *vetor){
	acquire();
	int i = indice_temp;
	indice_temp += 1;
	release();

	int soma_parcial = 0;

	for(int j=i*parcelas; j<i*parcelas+parcelas; j++){
		soma_parcial += vetor[j];
	}

	acquire();
	acumulador += soma_parcial;
	release();

	pthread_exit(0);
}

void *func2(int8_t *vetor){
	acquire();
	int i = indice_temp;
	indice_temp += 1;
	release();

	int soma_parcial = 0;

	for(int j=i*parcelas; j<(i*parcelas+parcelas)+resto; j++){
		soma_parcial += vetor[j];
	}

	acquire();
	acumulador += soma_parcial;
	release();

	pthread_exit(0);
}

void acquire(){
    while (__sync_lock_test_and_set(&lock, 1)) {}
}

void release(){
    lock = 0;
}

int main(int argc, char *argv[]){

	if(argc != 3 || atoi(argv[1])<=0 || atoi(argv[2])<=0) {
		printf("Formato: ./spinlock vector_size number_of_threads \n");
		return 0;
	}

	int8_t* vector = calloc(atoi(argv[1]), 1);

	parcelas = atoi(argv[1])/atoi(argv[2]);
	resto = atoi(argv[1])%atoi(argv[2]);

	srand(time(NULL));

	for(int i=0; i<atol(argv[1]); i++) {
		vector[i] = random_value();
	}

	pthread_t * thread = malloc(sizeof(pthread_t)*atoi(argv[2]));
    
	struct timespec start, finish;
	double elapsed;
	clock_gettime(CLOCK_MONOTONIC, &start);

	for(int i=0; i<atoi(argv[2]); i++) {
		if(i != atoi(argv[2]) - 1) {
			pthread_create(&thread[i], NULL, (void *(*)(void *)) func1, vector);
		}
		else {
			pthread_create(&thread[i], NULL, (void *(*)(void *)) func2, vector);
		}   
    }

	for(int i=0; i<atoi(argv[2]); i++) {
		pthread_join(thread[i], NULL);
	}

	clock_gettime(CLOCK_MONOTONIC, &finish);
	elapsed = (finish.tv_sec - start.tv_sec);
	elapsed += (finish.tv_nsec - start.tv_nsec) / 1000000000.0;

	printf("Valor total do acumulador: %d\n", acumulador);
	printf("Tempo total para execução da soma: %f\n", elapsed);
	return(0);
}