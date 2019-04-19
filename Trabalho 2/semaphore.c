#include <pthread.h>
#include <stdio.h>
#include <stdlib.h> 
#include <time.h>
#include <semaphore.h>
#include <unistd.h>

// gcc semaphore.c -pthread -o semaphore

#define last_number 100000 // Define quantos números serão examinados antes do programa terminar sua execução
#define vector_size 16 // Define o tamanho do buffer

volatile int numeros_examinados = 0;
int vector[vector_size] = {0};

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
	for(int i=0; i<vector_size; i++) {
		if(vetor[i] == 0) {
			return i;
		}
	}
	printf("Erro no produtor");
	return -1;
}

int find_occupied_space(int *vetor){
	for(int i=0; i<vector_size; i++) {
		if(vetor[i] != 0) {
			return i;
		}
	}
	printf("Erro no consumidor");
	return -1;
}

int random_value(){
	int valor = (rand() % 10000000) + 1;
	return valor;
}

void *produtor(int *vetor){
	while(1) {

		// Inicio da região crítica
		sem_wait(&empty);
		sem_wait(&mutex);

		int posicao = find_free_space(vetor);
		if(numeros_examinados < last_number) {
			vetor[posicao] = random_value();
		}

		else {
			// Finaliza a execução da thread, e aumenta os contadores de todos
			// os semáforos relacionados ao vetor para evitar deadlock
			sem_post(&mutex);
			sem_post(&empty);
			sem_post(&full);
			pthread_exit(0);
		}

		// Fim da região crítica
		sem_post(&mutex);
		sem_post(&full);
	}
}

void *consumidor(int *vetor){
	while(1) {

		// Início da região crítica
		sem_wait(&full);
		sem_wait(&mutex);

		if(numeros_examinados >= last_number){
			// Finaliza a execução da thread, e aumenta os contadores de todos
			// os semáforos relacionados ao vetor para evitar deadlock
			sem_post(&mutex);
			sem_post(&full);
			sem_post(&empty);
			pthread_exit(0);
		}

		int posicao = find_occupied_space(vetor);

		int valor = vetor[posicao];
		vetor[posicao] = 0;
		numeros_examinados += 1;

		// Fim da região crítica
		sem_post(&mutex);
		sem_post(&empty);

		int resultado = is_prime(valor);

		if(resultado == 1) {
			sem_wait(&primos);
			valores_primos += 1;
			sem_post(&primos);
		}
		else {
			sem_wait(&nao_primos);
			valores_nao_primos += 1;
			sem_post(&nao_primos);
		}
	}
}

int main(int argc, char *argv[]){

	if(argc != 3 || atoi(argv[1])<=0 || atoi(argv[2])<=0) {
		printf("Formato: ./semaphore threads_produtor threads_consumidor \n");
		return 0;
	}

	pthread_t * produtor_thread = malloc(sizeof(pthread_t)*atoi(argv[1]));
	pthread_t * consumidor_thread = malloc(sizeof(pthread_t)*atoi(argv[2]));

	sem_init(&mutex, 0, 1);
	sem_init(&empty, 0, vector_size);
	sem_init(&full, 0, 0);

	sem_init(&primos, 0, 1);
	sem_init(&nao_primos, 0, 1);

	srand(time(NULL));

	struct timespec start, finish;
	double elapsed;
	clock_gettime(CLOCK_MONOTONIC, &start);

	for(int i=0; i<atoi(argv[1]); i++) {
    		pthread_create(&produtor_thread[i], NULL, (void *(*)(void *)) produtor, vector);  
	}

	for(int i=0; i<atoi(argv[2]); i++) {
    		pthread_create(&consumidor_thread[i], NULL, (void *(*)(void *)) consumidor, vector);  
	}

	for(int i=0; i<atoi(argv[2]); i++) {
    		pthread_join(consumidor_thread[i], NULL);
	}

	clock_gettime(CLOCK_MONOTONIC, &finish);
	elapsed = (finish.tv_sec - start.tv_sec);
	elapsed += (finish.tv_nsec - start.tv_nsec) / 1000000000.0;

	printf("Execução completa.\n");
	printf("Valores primos: %d\n", valores_primos);
	printf("Valores não primos: %d\n", valores_nao_primos);
	printf("Tempo total para consumir 100000 números: %f\n", elapsed);

	return(0);
}