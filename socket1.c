#include <unistd.h> 
#include <stdio.h> 
#include <sys/socket.h> 
#include <stdlib.h> 
#include <netinet/in.h> 
#include <string.h>
#include <time.h>

#define PORT 8080 
#define SIZE 21

int main(int argc, char *argv[]) 
{ 

    if (argc != 2) {
        printf("Formato: ./socket1 number_of_values \n");
        return 0;
    }

    int server_fd, new_socket, valread; 
    struct sockaddr_in address; 
    int opt = 1; 
    int addrlen = sizeof(address);
    srand(time(NULL));

    char valor_enviado[SIZE];
    int valor = 1;
    int contador = 0;
    char primo_recebido[SIZE];
    int primo;
       
    server_fd = socket(AF_INET, SOCK_STREAM, 0); 
    setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR | SO_REUSEPORT, &opt, sizeof(opt));

    address.sin_family = AF_INET; 
    address.sin_addr.s_addr = INADDR_ANY; 
    address.sin_port = htons( PORT ); 
       
    bind(server_fd, (struct sockaddr *)&address, sizeof(address));
    listen(server_fd, 3);

    new_socket = accept(server_fd, (struct sockaddr *)&address, (socklen_t*)&addrlen);

    while (contador < atoi(argv[1])) {
        int valor_aleatorio = (rand() % 99) + 1;
        valor += valor_aleatorio;

        sprintf(valor_enviado, "%d", valor);
        send(new_socket , valor_enviado, SIZE, 0 );

        valread = read( new_socket, primo_recebido, SIZE);
        primo = atoi(primo_recebido);

        if (primo == 1) {
            printf("%d é primo.\n", valor);
        }

        else {
            printf("%d não é primo.\n", valor);
        }

        contador++;
    }

    send(new_socket, "0", SIZE, 0 );
    exit(0);
    return 0; 
}