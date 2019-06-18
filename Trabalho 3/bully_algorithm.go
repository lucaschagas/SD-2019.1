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
	address = "239.0.0.0"
	maxDatagramSize = 256
	multicastPort = 9999
)

var comando_buffer bytes.Buffer

var mensagens_enviadas[5] int 
var mensagens_recebidas[5] int

var lider_atual int
var lider_status int
var process_state int
var process_id int
var process_rank int
var ongoing_election int
var teste_ok int

var comando_buffer_lock = &sync.Mutex{}
var mensagens_enviadas_lock = &sync.Mutex{}
var mensagens_recebidas_lock = &sync.Mutex{}
var lider_atual_lock = &sync.Mutex{}
var lider_status_lock = &sync.Mutex{}
var process_state_lock = &sync.Mutex{}
var ongoing_election_lock = &sync.Mutex{}
var ok_lock = &sync.Mutex{}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	s := string(b)
	comando_buffer_lock.Lock()
	comando_buffer.WriteString(s)
	comando_buffer_lock.Unlock()
}

func send_message(message string, port int) {
	conn := NewBroadcaster(address + ":" + strconv.Itoa(port))
	conn.Write([]byte(message))
}

func NewBroadcaster(address string) (*net.UDPConn) {
	addr, _ := net.ResolveUDPAddr("udp", address)
	conn, _ := net.DialUDP("udp", nil, addr)
	return conn
}

func Listen(address string, handler func(*net.UDPAddr, int, []byte), port int) {
	addr, _ := net.ResolveUDPAddr("udp", address + ":" + strconv.Itoa(port))
	conn, _ := net.ListenMulticastUDP("udp", nil, addr)
	conn.SetReadBuffer(maxDatagramSize)

	for {
		buffer := make([]byte, maxDatagramSize)
		numBytes, src, _ := conn.ReadFromUDP(buffer)

		handler(src, numBytes, buffer)
	}
}

func start_election() {
	ongoing_election_lock.Lock()
	ongoing_election = 1
	ongoing_election_lock.Unlock()

	ok_lock.Lock()
	teste_ok = 0
	ok_lock.Unlock()

	send_message("0" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", multicastPort)
	update_sent_messages(0)
	time.Sleep(6 * time.Second)

	ok_lock.Lock()
	if teste_ok == 0 {
		ok_lock.Unlock()
		send_message("2" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", multicastPort)
		update_sent_messages(2)
	} else {
		ok_lock.Unlock()
	}

	ongoing_election_lock.Lock()
	ongoing_election = 0
	ongoing_election_lock.Unlock()
}

func check_leader_alive() (answer int) {
	lider_status_lock.Lock()
	lider_status = 0
	lider_status_lock.Unlock()

	lider_atual_lock.Lock()
	send_message("3" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", lider_atual)
	lider_atual_lock.Unlock()

	update_sent_messages(3)

	time.Sleep(5 * time.Second)
	lider_status_lock.Lock()
	if lider_status == 1 {
		lider_status_lock.Unlock()
		return 1
	} else {
		lider_status_lock.Unlock()
		return 0
	}
}

func update_received_messages(i int) {
	mensagens_recebidas_lock.Lock()
	mensagens_recebidas[i] += 1
	mensagens_recebidas_lock.Unlock()
}

func update_sent_messages(i int) {
	mensagens_enviadas_lock.Lock()
	mensagens_enviadas[i] += 1
	mensagens_enviadas_lock.Unlock()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	process_id = rand.Intn(7998)+1125
	process_rank = process_id

	process_state = 1 //0: Falha; 1: Funcionando
	lider_status = 1 //0: Morto; 1: Vivo
	ongoing_election = 0 //0: Sem eleição em execução; 1: Eleição em execução
	teste_ok = 0 //0: Nenhum OK recebido durante eleição; 1: Ao menos 1 OK recebido durante eleição

	fmt.Println("ID do processo:", process_id)

	var wg sync.WaitGroup
	wg.Add(1)

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
			  		fmt.Println("alive - Verifica se o líder atual está vivo; caso não esteja, avisa todos os processos que o líder está morto.")
			  		fmt.Println("clear - Reseta a quantidade de mensagens enviadas e recebidas imprimida pelo comando stats.")

			  	case "election":

			  		process_state_lock.Lock()
			  		if process_state == 0 {
			  			process_state_lock.Unlock()
			  			fmt.Println("Processo está simulando uma falha, e não pode enviar mensagens.")
			  		} else {
			  			process_state_lock.Unlock()		
			  		}

			  		process_state_lock.Lock()
			  		if process_state == 1 {
			  			process_state_lock.Unlock()

			  			ongoing_election_lock.Lock()
						if ongoing_election == 0 {
							ongoing_election_lock.Unlock()
							go start_election()
						} else {
							ongoing_election_lock.Unlock()
						}

					} else {
			  			process_state_lock.Unlock()		
			  		}

			  	case "fail":

			  		process_state_lock.Lock()
			  		if process_state == 0 {
			  			process_state_lock.Unlock()
			  			fmt.Println("O processo já está simulando uma falha.")
			  		} else {
			  			process_state_lock.Unlock()		
			  		}

			  		process_state_lock.Lock()
			  		if process_state == 1 {
			  			process_state = 0
			  			process_state_lock.Unlock()
			  			fmt.Println("O processo agora está simulando uma falha.")
			  		} else {
			  			process_state_lock.Unlock()		
			  		}

			  	case "recover":

			  		process_state_lock.Lock()
			  		if process_state == 1 {
			  			process_state_lock.Unlock()
			  			fmt.Println("O processo já está funcionando normalmente.")
			  		} else {
			  			process_state_lock.Unlock()		
			  		}

			  		process_state_lock.Lock()
			  		if process_state == 0 {
			  			process_state = 1
			  			process_state_lock.Unlock()
			  			fmt.Println("O processo está agora funcionando normalmente.")
			  		} else {
			  			process_state_lock.Unlock()		
			  		}

			  	case "stats":
			  		lider_atual_lock.Lock()
			  		fmt.Printf("Líder atual do sistema: %d\n", lider_atual)
			  		lider_atual_lock.Unlock()
			  		mensagens_enviadas_lock.Lock()
			  		fmt.Printf("Mensagens enviadas- ELEICAO: %d | OK: %d | LIDER: %d | VIVO: %d | VIVO_OK: %d\n",
			  			mensagens_enviadas[0], mensagens_enviadas[1], mensagens_enviadas[2], mensagens_enviadas[3], mensagens_enviadas[4])
			  		mensagens_enviadas_lock.Unlock()
			  		mensagens_recebidas_lock.Lock()
			  		fmt.Printf("Mensagens recebidas- ELEICAO: %d | OK: %d | LIDER: %d | VIVO: %d | VIVO_OK: %d\n",
			  			mensagens_recebidas[0], mensagens_recebidas[1], mensagens_recebidas[2], mensagens_recebidas[3], mensagens_recebidas[4])
			  		mensagens_recebidas_lock.Unlock()

			  	case "alive":

			  		process_state_lock.Lock()
			  		if process_state == 0 {
			  			process_state_lock.Unlock()
			  			fmt.Println("Processo está simulando uma falha, e não pode enviar mensagens.")
			  		} else {
			  			process_state_lock.Unlock()
			  			manual_check := check_leader_alive()

			  			if manual_check == 1 {
							fmt.Println("Líder está vivo.")
			  			} else {
			  				fmt.Println("Líder está morto.")
			  			}		
			  		}

			  	case "clear":

					mensagens_recebidas_lock.Lock()
					mensagens_enviadas_lock.Lock()
					for i:= 0; i < 5; i++ {
						mensagens_enviadas[i] = 0
						mensagens_recebidas[i] = 0
					}
					mensagens_recebidas_lock.Unlock()
					mensagens_enviadas_lock.Unlock()
		  	}
	  	}
	}()

	// Leitura de mensagens enviadas para todos os processos
	go Listen(address, msgHandler, multicastPort)

	// Leitura de mensagens para este processo
	go Listen(address, msgHandler, process_id)

	// Processamento de mensagens
	go func() {

		for true {
			comando_buffer_lock.Lock()
			comando_atual_bytes, _ := comando_buffer.ReadBytes(';')
			comando_buffer_lock.Unlock()
			comando_atual_string := string(bytes.Trim(comando_atual_bytes, "\x00"))

			if comando_atual_string != "" {

				string_split := strings.Split(comando_atual_string, "|")
				time.Sleep(30 * time.Millisecond)
				codigo, processo_string, valor := string_split[0], string_split[1], string_split[2]

				valor_temp := strings.TrimSuffix(valor, ";")
				valor_int, _ := strconv.Atoi(valor_temp)
				codigo_int, _ := strconv.Atoi(codigo)
				processo_int, _ := strconv.Atoi(processo_string)

				/*if processo_int == process_id && lider_atual != process_id {
					process_id = rand.Intn(7998)+1125
				}*/

				//fmt.Println("Código:", codigo, "| Valor:", valor_int)

				switch codigo_int {
					case 0:
						//fmt.Println("Código ELEIÇÃO")
						update_received_messages(0)

						process_state_lock.Lock()
						if process_state == 1 {
							process_state_lock.Unlock()

							if process_rank > valor_int {

								send_message("1" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", processo_int)
								update_sent_messages(1)

								ongoing_election_lock.Lock()
								if ongoing_election == 0 {
									ongoing_election_lock.Unlock()
									go start_election()
								} else {
									ongoing_election_lock.Unlock()
								}
							}
						} else {
			  				process_state_lock.Unlock()		
			  			}

					case 1:
						//fmt.Println("Código OK")
						update_received_messages(1)

						ok_lock.Lock()
						teste_ok = 1
						ok_lock.Unlock()

					case 2:
						//fmt.Println("Código LÍDER")
						update_received_messages(2)

						lider_atual_lock.Lock()
						lider_atual = valor_int
						fmt.Println("Líder atual:", lider_atual)
						lider_atual_lock.Unlock()

					case 3:
						//fmt.Println("Código VIVO")
						update_received_messages(3)

						process_state_lock.Lock()
						if process_state == 1 {
							process_state_lock.Unlock()
							send_message("4" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", processo_int)
							update_sent_messages(4) 
						} else {
			  				process_state_lock.Unlock()		
			  			}

					case 4:
						//fmt.Println("Código VIVO_OK")
						update_received_messages(4)

						lider_status_lock.Lock()
						lider_status = 1
						lider_status_lock.Unlock()
				}
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// Detecção automática do líder
	go func() {
		for true {
			time.Sleep(400 * time.Second)

			if process_state == 1 {
				auto_check := check_leader_alive()

				if auto_check == 0 {
					ongoing_election_lock.Lock()

					if ongoing_election == 0 {
						ongoing_election_lock.Unlock()
						go start_election()
					} else {
						ongoing_election_lock.Unlock()
					}
			  	}
			}
	  	}
	}()

	wg.Wait()
}