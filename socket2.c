#include <unistd.h> 
#include <stdio.h> 
#include <sys/socket.h> 
#include <stdlib.h> 
#include <netinet/in.h> 
#include <string.h> 
#include <arpa/inet.h>
#include <time.h>

#define PORT 8080 
#define SIZE 21

int is_prime(int n) {

    int i, flag = 0;

    for(i = 2; i <= n/2; ++i) {
        if (n%i == 0) {
            flag = 1;
            break;
        }
    }

    if (n == 1) {
        printf("1 não é primo.\n");
        return 0;
    }

    else {
        if (flag == 0) {
            printf("%d é primo.\n", n);
            return 1;
        }
        else {
            printf("%d não é primo.\n", n);
            return 0;
        }
    }
}
   
int main(int argc, char *argv[]) 
{ 

    int sock = 0, valread; 
    struct sockaddr_in serv_addr;

    char valor_recebido[SIZE];
    int valor_final;
    int primo;
    char primo_enviado[SIZE];

    sock = socket(AF_INET, SOCK_STREAM, 0);
    memset(&serv_addr, '0', sizeof(serv_addr)); 
   
    serv_addr.sin_family = AF_INET; 
    serv_addr.sin_port = htons( PORT ); 
       
    inet_pton(AF_INET, "127.0.0.1", &serv_addr.sin_addr);
    connect(sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr));

    while (valor_final != 0) {
        valread = read( sock, valor_recebido, SIZE);
        valor_final = atoi(valor_recebido);

        if (valor_final > 0) {
            primo = is_prime(valor_final);
        }

        sprintf(primo_enviado, "%d", primo);
        send(sock, primo_enviado, SIZE , 0 );
    }

    exit(0); 
    return 0; 
} 