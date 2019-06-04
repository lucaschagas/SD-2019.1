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
	address = "239.0.0.0:9999"
	maxDatagramSize = 256
)

var comando_buffer bytes.Buffer
var comando_atual string

var mensagens_enviadas[5] int 
var mensagens_recebidas[5] int

var lider_atual int
var lider_status int
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

func send_message(message string) {
	conn := NewBroadcaster(address)
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
	process_state = 1 //0: Falha; 1: Funcionando
	lider_status = 1 //0: Morto; 1: Vivo

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
			  		fmt.Println("fail - Emula uma falha de processo, de modo que mensagens recebidas serão ignoradas. Use o comando \"recover\" para retomar o processo ao funcionamento normal; o líder atual ainda será atualizado.")
			  		fmt.Println("recover - Recupera um processo que esteja emulando falha.")
			  		fmt.Println("stats - Imprime o atual líder do sistema, e a quantidade de mensagens enviadas e recebidas de cada tipo.")
			  		fmt.Println("alive - Verifica se o líder atual está vivo.")

			  	case "election":

			  		if process_state == 0 {
			  			fmt.Println("Processo está simulando uma falha, e não pode enviar mensagens.")
			  		}

			  		if process_state == 1 {

				  		election_status_start = 0 //0: Não recebeu OK; 1: Recebeu OK
				  		send_message("0" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
				  		time.Sleep(time.Duration(3) * time.Second)

				  		if election_status_start == 0 {
					  		send_message("2" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
					  		fmt.Println("Comando executado com sucesso.")
					  	} else {
					  		fmt.Println("Existe processo ativo com rank maior·")
					  	}
					}

			  	case "fail":

			  		if process_state == 0 {
			  			fmt.Println("O processo já está simulando uma falha.")
			  		}

			  		if process_state == 1 {
			  			process_state = 0
			  			fmt.Println("O processo agora está simulando uma falha.")
			  		}

			  	case "recover":

			  		if process_state == 0 {
			  			fmt.Println("O processo já está funcionando normalmente.")
			  		}

			  		if process_state == 0 {
			  			process_state = 1
			  			fmt.Println("O processo está agora funcionando normalmente.")
			  		}

			  	case "stats":
			  		fmt.Printf("Mensagens enviadas- ELEICAO: %d | OK: %d | LIDER: %d | VIVO: %d | VIVO_OK: %d\n", mensagens_enviadas[0], mensagens_enviadas[1], mensagens_enviadas[2], mensagens_enviadas[3], mensagens_enviadas[4])
			  		fmt.Printf("Mensagens recebidas- ELEICAO: %d | OK: %d | LIDER: %d | VIVO: %d | VIVO_OK: %d\n", mensagens_recebidas[0], mensagens_recebidas[1], mensagens_recebidas[2], mensagens_recebidas[3], mensagens_recebidas[4])
			  		fmt.Println("Comando executado com sucesso.")

			  	case "alive":

			  		if process_state == 0 {
			  			fmt.Println("Processo está simulando uma falha, e não pode enviar mensagens.")
			  		}

			  		if process_state == 1 {

				  		send_message("3" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")

				  		time.Sleep(time.Duration(3) * time.Second)
				  		if lider_status == 1 {
				  			fmt.Println("Líder está vivo.")
				  		} else {
				  			time.Sleep(time.Duration(7) * time.Second)
				  			if lider_status == 1 {
				  				fmt.Println("Líder está vivo.")
				  			} else {
				  				fmt.Println("Líder está morto.")
				  				lider_status = 0
				  			}
				  		}
				  	}
		  	}
	  	}
	}()

	// Leitura de mensagens na rede
	go Listen(address, msgHandler)

	// Processamento de mensagens
	go func() {

		for true {
			comando_atual_bytes, _ := comando_buffer.ReadBytes(';')
			comando_atual_string := string(bytes.Trim(comando_atual_bytes, "\x00"))

			if comando_atual_string != "" {

				string_split := strings.Split(comando_atual_string, "|")
				time.Sleep(time.Duration(1) * time.Second)
				codigo, _, valor := string_split[0], string_split[1], string_split[2]

				valor_temp := strings.TrimSuffix(valor, ";")
				valor_int, _ := strconv.Atoi(valor_temp)
				codigo_int, err := strconv.Atoi(codigo)

				if err != nil {
					fmt.Println(err.Error())
				}

				fmt.Println("Código:", codigo, "| Valor:", valor_int)

				switch codigo_int {
					case 0:
						fmt.Println("Código ELEIÇÃO")

						if process_state == 1 {

							if process_rank > int32(valor_int) {

								send_message("1" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")

								election_status_receive = 0
								send_message("0" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
								time.Sleep(time.Duration(3) * time.Second)
								
								if election_status_receive == 0 {
									send_message("2" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
								}
							}
						}

					case 1:
						fmt.Println("Código OK")
						election_status_start = 1
						election_status_receive = 1

					case 2:
						fmt.Println("Código LÍDER")
						if lider_atual < valor_int && lider_status == 1{
							lider_atual = valor_int
						} else if lider_status == 0 {
							lider_atual = valor_int
						}

					case 3:
						fmt.Println("Código VIVO")

						if process_state == 1 {
							if process_rank == int32(lider_atual) {
								send_message("4" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";")
							}
						}

					case 4:
						fmt.Println("Código VIVO_OK")
						new_election_status = 1
						lider_status = 1

					/*
					default:
						fmt.Println("DEFAULT")
						fmt.Println("Código:", codigo, "| Valor:", valor_int)
						goto myswitch
					*/
				}
			} else {
				time.Sleep(time.Duration(10) * time.Millisecond)
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
	  		fmt.Println("Líder atual:", lider_atual)
	  	}
	}()

	wg.Wait()

}