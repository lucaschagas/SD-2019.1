#include <pthread.h>
#include <stdio.h>
#include <stdlib.h> 
#include <time.h> 

volatile int lock = 0;
volatile int acumulador = 0;
volatile int indice_temp = 0;
volatile int parcelas = 0;
volatile int resto = 0;

void acquire();

void release();

int random_value(){
	int valor = (rand() % 201) - 100;
	return valor;
}

void *func1(int *vetor) {
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

void *func2(int *vetor) {
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

void acquire() {
    while (__sync_lock_test_and_set(&lock, 1)) {}
}

void release() {
    __sync_synchronize();
    lock = 0;
}

int main(int argc, char *argv[]){

	if(argc != 3 || atol(argv[1])<=0 || atoi(argv[2])<=0) {
		printf("Formato: ./spinlock vector_size number_of_threads \n");
		return 0;
	}

	long temp = atol(argv[1]);
	printf("%ld\n", temp);

	int vector[temp];

	int temp_soma = 0;

	parcelas = atol(argv[1])/atoi(argv[2]);
	resto = atol(argv[1])%atoi(argv[2]);

	srand(time(NULL));

	for(int i=0; i<atol(argv[1]); i++) {
		vector[i] = random_value();
		temp_soma += vector[i];
	}

    pthread_t * thread = malloc(sizeof(pthread_t)*atoi(argv[2]));
    
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

    printf("Valor correto da soma: %d\n", temp_soma);
    printf("Valor total do acumulador: %d\n", acumulador);
	return(0);
}