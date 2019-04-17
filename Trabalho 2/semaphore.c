#include <pthread.h>
#include <stdio.h>
#include <stdlib.h> 
#include <time.h>
#include <semaphore.h>
#include <unistd.h>
#include <string.h>

// gcc semaphore.c -lpthread -lrt -o semaphore

#define last_number 100000 // Define quantos números serão examinados antes do programa terminar sua execução

volatile int numeros_examinados = 0;

sem_t mutex;
sem_t empty;
sem_t full;

sem_t primos;
sem_t nao_primos;

int valores_primos = 0;
int valores_nao_primos = 0;

int is_prime(int n){
	int i, flag = 0;

	for(i = 2; i <= n/2; ++i) {
		if (n%i == 0) {
			flag = 1;
			break;
		}
	}

	if (n == 1) {
		return 0; // Não é primo
	}

	else {
		if (flag == 0) {
			return 1; // É primo
		}
		else {
			return 0; // Não é primo
        	}
	}
}

int find_free_space(int *vetor){
	for(int i=0; i<sizeof(vetor); i++) {
		if(vetor[i] == 0) {
			return i;
		}
	}
	return -1;
}

int find_occupied_space(int *vetor){
	for(int i=0; i<sizeof(vetor); i++) {
		if(vetor[i] != 0) {
			return i;
		}
	}
	return -1;
}

int random_value(){
	int valor = (rand() % 10000000) + 1;
	return valor;
}

void *produtor(int *vetor){
	while(numeros_examinados < last_number) {

		// Inicio da região crítica
		sem_wait(&empty);
		sem_wait(&mutex);

		int posicao = find_free_space(vetor);
		if(numeros_examinados < last_number) {
			vetor[posicao] = random_value();
		}

		// Fim da região crítica
		sem_post(&mutex);
		sem_post(&full);
	}

	pthread_exit(0);
}

void *consumidor(int *vetor){
	while(numeros_examinados < last_number) {

		// Início da região crítica
		sem_wait(&full);
		sem_wait(&mutex);

		int posicao = find_occupied_space(vetor);

		int valor = vetor[posicao];
		vetor[posicao] = 0;
		numeros_examinados += 1;

		// Fim da região crítica
		sem_post(&mutex);
		sem_post(&empty);

		if(numeros_examinados == last_number + 1) {
			pthread_exit(0);
		}

		int resultado = is_prime(valor);

		if(resultado == 1) {
			// printf("Valor %d é primo\n", valor);
			sem_wait(&primos);
			valores_primos += 1;
			sem_post(&primos);
		}
		else {
			// printf("Valor %d não é primo\n", valor);
			sem_wait(&nao_primos);
			valores_nao_primos += 1;
			sem_post(&nao_primos);
		}
	}

	pthread_exit(0);
}

int main(int argc, char *argv[]){

	if(argc != 4 || atoi(argv[1])<=0 || atoi(argv[2])<=0 || atoi(argv[3])<=0) {
		printf("Formato: ./semaphore threads_produtor threads_consumidor buffer_size \n");
		return 0;
	}

	int* vector = calloc(atoi(argv[3]), sizeof(int));

	pthread_t * produtor_thread = malloc(sizeof(pthread_t)*atoi(argv[1]));
	pthread_t * consumidor_thread = malloc(sizeof(pthread_t)*atoi(argv[2]));

	sem_init(&mutex, 0, 1);
	sem_init(&empty, 0, atoi(argv[3]));
	sem_init(&full, 0, 0);

	sem_init(&primos, 0, 1);
	sem_init(&nao_primos, 0, 1);

	srand(time(NULL));

	for(int i=0; i<atoi(argv[1]); i++) {
    		pthread_create(&produtor_thread[i], NULL, (void *(*)(void *)) produtor, vector);  
	}

	for(int i=0; i<atoi(argv[2]); i++) {
    		pthread_create(&consumidor_thread[i], NULL, (void *(*)(void *)) consumidor, vector);  
	}

	for(int i=0; i<atoi(argv[2]); i++) {
    		pthread_join(consumidor_thread[i], NULL);
	}

	printf("Execução completa.\n");
	printf("Valores primos: %d\n", valores_primos);
	printf("Valores não primos: %d\n", valores_nao_primos);

	return(0);
}
