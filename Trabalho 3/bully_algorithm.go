package main

import (
	"fmt"
	"math/rand"
	"time"
	"sync"
	"net"
	"bytes"
	"strings"
	"strconv"
)

const (
	address_all = "239.0.0.0:9999"
	address_direct = "239.0.0.1:9999"
	maxDatagramSize = 8192
)

var comando_buffer bytes.Buffer
var comando_atual string

var mensagens_enviadas[5] int 
var mensagens_recebidas[5] int

var lider_atual int
var process_id int32
var process_rank int32
var process_state int
var election_status_start int
var election_status_receive int
var new_election_status int

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	s := string(b)
	comando_buffer.WriteString(s)
}

func send_message_all(message string) {
	conn := NewBroadcaster(address_all)
	conn.Write([]byte(message))
}

func send_message_direct(message string) {
	conn := NewBroadcaster(address_direct)
	conn.Write([]byte(message))
}

func NewBroadcaster(address string) (*net.UDPConn) {
	addr, _ := net.ResolveUDPAddr("udp", address)
	conn, _ := net.DialUDP("udp", nil, addr)
	return conn
}

func Listen(address string, handler func(*net.UDPAddr, int, []byte)) {
	addr, _ := net.ResolveUDPAddr("udp", address)
	conn, _ := net.ListenMulticastUDP("udp", nil, addr)
	conn.SetReadBuffer(maxDatagramSize)

	for {
		buffer := make([]byte, maxDatagramSize)
		numBytes, src, _ := conn.ReadFromUDP(buffer)

		handler(src, numBytes, buffer)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	process_id = rand.Int31()
	process_rank = process_id
	process_state = 0 //0: Funcionando; 1: Falha

	var wg sync.WaitGroup
	wg.Add(3)

	// Interface com o usuário
	go func() {

		var comando string
		for true {

		  	fmt.Println("Digite o comando desejado; para mais informações, digite \"help\".")
		  	fmt.Scanln(&comando)

			switch comando {
		  		case "help":
			  		fmt.Println("Comandos disponíveis:")
			  		fmt.Println("election - Verifica se o líder atual está operacional; caso não esteja, inicia uma eleição.")
			  		fmt.Println("fail - Emula uma falha de processo, de modo que mensagens recebidas serão ignoradas; o líder atual ainda será atualizado.")
			  		fmt.Println("recover - Recupera um processo que esteja emulando falha.")
			  		fmt.Println("stats - Imprime o atual líder do sistema, e a quantidade de mensagens enviadas e recebidas de cada tipo.")

			  	case "election":
			  		election_status_start = 0 //0: Não recebeu OK; 1: Recebeu OK
			  		send_message_all("0" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
			  		time.Sleep(time.Duration(5) * time.Second)

			  		if election_status_start == 0 {
				  		send_message_all("2" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
				  		fmt.Println("Comando executado com sucesso.")
				  	} else {
				  		fmt.Println("Existe processo ativo com rank maior·")
				  	}

			  	case "fail":
			  		process_state = 1
			  		fmt.Println("Comando executado com sucesso.")

			  	case "recover":
			  		process_state = 0
			  		fmt.Println("Comando executado com sucesso.")

			  	case "stats":
			  		fmt.Printf("Mensagens enviadas- ELEICAO: %d | OK: %d | LIDER: %d | VIVO: %d | VIVO_OK: %d\n", mensagens_enviadas[0], mensagens_enviadas[1], mensagens_enviadas[2], mensagens_enviadas[3], mensagens_enviadas[4])
			  		fmt.Printf("Mensagens recebidas- ELEICAO: %d | OK: %d | LIDER: %d | VIVO: %d | VIVO_OK: %d\n", mensagens_recebidas[0], mensagens_recebidas[1], mensagens_recebidas[2], mensagens_recebidas[3], mensagens_recebidas[4])
			  		fmt.Println("Comando executado com sucesso.")
		  	}
	  	}
	}()

	// Processamento de mensagens
	go func() {

		go Listen(address_all, msgHandler)

		for true {
			comando_atual, erro_buffer := comando_buffer.ReadString(';')

			var erro_atual string = ""

			if erro_buffer != nil {
				if erro_buffer.Error() == "EOF" {
					erro_atual = erro_buffer.Error()
				}
			}

			if comando_atual != "" && erro_atual != "EOF" {
				erro_atual = ""

				string_split := strings.Split(comando_atual, "|")
				time.Sleep(time.Duration(1) * time.Second)
				codigo, _, valor := string_split[0], string_split[1], string_split[2]

				valor_temp := strings.TrimSuffix(valor, ";")
				valor_int, _ := strconv.Atoi(valor_temp)
				codigo_temp, _ := strconv.Atoi(codigo)

				fmt.Println(codigo, valor, lider_atual, codigo_temp)

				switch codigo_temp {
					case 0:
						fmt.Println("Código ELEIÇÃO")
						if process_rank > int32(valor_int) {

							send_message_all("1" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
							time.Sleep(time.Duration(1) * time.Second)

							election_status_receive = 0
							send_message_all("0" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
							time.Sleep(time.Duration(3) * time.Second)
							
							if election_status_receive == 0 {
								send_message_all("2" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
							}
						}

					case 1:
						fmt.Println("Código OK")
						election_status_start = 1
						election_status_receive = 1

					case 2:
						lider_atual = valor_int

					case 3:
						if process_rank == int32(lider_atual) {
							send_message_direct("4" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
						}

					case 4:
						new_election_status = 1

					default:
						fmt.Println("DEFAULT")
						fmt.Println(codigo)
				}
			} else {
				time.Sleep(time.Duration(1) * time.Second)
			}
		}
	}()

	// Detecção do líder
	go func() {

		//go Listen(address_direct, msgHandler)

		for true {

			new_election_status = 0
			time.Sleep(time.Duration(5) * time.Second)
/*
			if process_rank != int32(lider_atual) {

				send_message_all("3" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
		  		time.Sleep(time.Duration(3) * time.Second)

		  		if new_election_status == 0 {
		  			send_message_all("0" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
		  		}
	  		}*/
	  		fmt.Println(lider_atual)
	  	}
	}()

	wg.Wait()

}