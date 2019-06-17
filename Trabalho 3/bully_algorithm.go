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

var mensagens_enviadas[6] int 
var mensagens_recebidas[6] int

var lider_atual int
var lider_status int
var process_state int
var election_status_start int
var election_status_receive int

var comando_buffer_lock = &sync.Mutex{}
var mensagens_enviadas_lock = &sync.Mutex{}
var mensagens_recebidas_lock = &sync.Mutex{}
var lider_atual_lock = &sync.Mutex{}
var lider_status_lock = &sync.Mutex{}
var process_state_lock = &sync.Mutex{}

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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	process_id := rand.Intn(9997)+1
	process_rank := process_id

	fmt.Println("")
	fmt.Println(process_id)
	fmt.Println("")

	// Implementar semáforo para estados do líder e do processo
	process_state = 1 //0: Falha; 1: Funcionando
	lider_status = 1 //0: Morto; 1: Vivo

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

				  		election_status_start = 0 //0: Não recebeu OK; 1: Recebeu OK
				  		send_message("0" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", multicastPort)
				  		mensagens_enviadas_lock.Lock()
				  		mensagens_enviadas[0] += 1
				  		mensagens_enviadas_lock.Unlock()
				  		time.Sleep(time.Duration(4) * time.Second)

				  		if election_status_start == 0 {
					  		send_message("2" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", multicastPort)
					  		mensagens_enviadas_lock.Lock()
					  		mensagens_enviadas[2] += 1
					  		mensagens_enviadas_lock.Unlock()
					  		fmt.Println("Comando executado com sucesso.")
					  	} else {
					  		fmt.Println("Existe processo ativo com rank maior·")
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
			  		mensagens_enviadas_lock.Lock()
			  		fmt.Printf("Mensagens enviadas- ELEICAO: %d | OK: %d | LIDER: %d | VIVO: %d | VIVO_OK: %d | LIDER_MORTO: %d\n", mensagens_enviadas[0], mensagens_enviadas[1], mensagens_enviadas[2], mensagens_enviadas[3], mensagens_enviadas[4], mensagens_enviadas[5])
			  		mensagens_enviadas_lock.Unlock()
			  		mensagens_recebidas_lock.Lock()
			  		fmt.Printf("Mensagens recebidas- ELEICAO: %d | OK: %d | LIDER: %d | VIVO: %d | VIVO_OK: %d | LIDER_MORTO: %d\n", mensagens_recebidas[0], mensagens_recebidas[1], mensagens_recebidas[2], mensagens_recebidas[3], mensagens_recebidas[4], mensagens_recebidas[5])
			  		mensagens_recebidas_lock.Unlock()
			  		fmt.Println("Comando executado com sucesso.")

			  	case "alive":

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

			  			lider_atual_lock.Lock()
				  		send_message("3" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", lider_atual)
				  		lider_atual_lock.Unlock()
				  		mensagens_enviadas_lock.Lock()
				  		mensagens_enviadas[3] += 1
				  		mensagens_enviadas_lock.Unlock()
				  		lider_status_lock.Lock()
				  		lider_status = 0
				  		lider_status_lock.Unlock()

				  		time.Sleep(time.Duration(3) * time.Second)
				  		lider_status_lock.Lock()
				  		if lider_status == 1 {
				  			lider_status_lock.Unlock()
				  			fmt.Println("Líder está vivo.")
				  		} else {
				  			lider_status_lock.Unlock()
				  			time.Sleep(time.Duration(3) * time.Second)
				  			lider_status_lock.Lock()
				  			if lider_status == 1 {
				  				lider_status_lock.Unlock()
				  				fmt.Println("Líder está vivo.")
				  			} else {
				  				lider_status_lock.Unlock()
				  				fmt.Println("Líder está morto.")
				  				send_message("5" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", multicastPort)
				  				mensagens_enviadas_lock.Lock()
				  				mensagens_enviadas[5] += 1
				  				mensagens_enviadas_lock.Unlock()
				  			}
				  		}
				  	} else {
			  			process_state_lock.Unlock()		
			  		}
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
				time.Sleep(time.Duration(1) * time.Second)
				codigo, processo_string, valor := string_split[0], string_split[1], string_split[2]

				valor_temp := strings.TrimSuffix(valor, ";")
				valor_int, _ := strconv.Atoi(valor_temp)
				codigo_int, _ := strconv.Atoi(codigo)
				processo_int, _ := strconv.Atoi(processo_string)

				fmt.Println("Código:", codigo, "| Valor:", valor_int)

				switch codigo_int {
					case 0:
						fmt.Println("Código ELEIÇÃO")
						mensagens_recebidas_lock.Lock()
						mensagens_recebidas[0] += 1
						mensagens_recebidas_lock.Unlock()

						process_state_lock.Lock()
						if process_state == 1 {
							process_state_lock.Unlock()

							if process_rank > valor_int {

								send_message("1" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", processo_int)
								mensagens_enviadas_lock.Lock()
								mensagens_enviadas[1] += 1
								mensagens_enviadas_lock.Unlock()

								election_status_receive = 0
								send_message("0" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", multicastPort)
								mensagens_enviadas_lock.Lock()
								mensagens_enviadas[0] += 1
								mensagens_enviadas_lock.Unlock()
								time.Sleep(time.Duration(4) * time.Second)
								
								if election_status_receive == 0 {
									send_message("2" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", multicastPort)
									mensagens_enviadas_lock.Lock()
									mensagens_enviadas[2] += 1
									mensagens_enviadas_lock.Unlock()
								}
							}
						} else {
			  				process_state_lock.Unlock()		
			  			}

					case 1:
						fmt.Println("Código OK")
						mensagens_recebidas_lock.Lock()
						mensagens_recebidas[1] += 1
						mensagens_recebidas_lock.Unlock()

						election_status_start = 1
						election_status_receive = 1

					case 2:
						fmt.Println("Código LÍDER")
						mensagens_recebidas_lock.Lock()
						mensagens_recebidas[2] += 1
						mensagens_recebidas_lock.Unlock()

						lider_atual_lock.Lock()
						lider_status_lock.Lock()
						if lider_atual < valor_int && lider_status == 1 {
							lider_atual = valor_int
						}

						if lider_status == 0 {
							lider_atual = valor_int
							lider_status = 1
						}
						lider_status_lock.Unlock()
						lider_atual_lock.Unlock()

					case 3:
						fmt.Println("Código VIVO")
						mensagens_recebidas_lock.Lock()
						mensagens_recebidas[3] += 1
						mensagens_recebidas_lock.Unlock()

						process_state_lock.Lock()
						if process_state == 1 {
							process_state_lock.Unlock()
							lider_atual_lock.Lock()
							if process_rank == lider_atual {
								lider_atual_lock.Unlock()
								send_message("4" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", processo_int)
								mensagens_enviadas_lock.Lock()
								mensagens_enviadas[4] += 1
								mensagens_enviadas_lock.Unlock()
							} else {
			  					lider_atual_lock.Unlock()		
			  				}
						} else {
			  				process_state_lock.Unlock()		
			  			}

					case 4:
						fmt.Println("Código VIVO_OK")
						mensagens_recebidas_lock.Lock()
						mensagens_recebidas[4] += 1
						mensagens_recebidas_lock.Unlock()

						lider_status_lock.Lock()
						lider_status = 1
						lider_status_lock.Unlock()

					case 5:
						fmt.Println("Código DEAD_LÍDER")
						mensagens_recebidas_lock.Lock()
						mensagens_recebidas[5] += 1
						mensagens_recebidas_lock.Unlock()
						
						lider_status_lock.Lock()
						lider_status = 0
						lider_status_lock.Unlock()
				}
			} else {
				time.Sleep(time.Duration(10) * time.Millisecond)
			}
		}
	}()

	// Detecção do líder
	go func() {

		for true {
			time.Sleep(time.Duration(40) * time.Second)

			lider_atual_lock.Lock()
			send_message("3" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", lider_atual)
			lider_temp := lider_atual
			lider_atual_lock.Unlock()
			mensagens_enviadas_lock.Lock()
			mensagens_enviadas[3] += 1
			mensagens_enviadas_lock.Unlock()
			lider_status_lock.Lock()
			lider_status = 0
			lider_status_lock.Unlock()

			time.Sleep(time.Duration(7) * time.Second)

			lider_status_lock.Lock()
			if lider_status == 0 && lider_temp == lider_atual {
				lider_status_lock.Unlock()
				send_message("5" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", multicastPort)
				mensagens_enviadas_lock.Lock()
				mensagens_enviadas[5] += 1
				mensagens_enviadas_lock.Unlock()
				
				election_status_receive = 0
				send_message("0" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", multicastPort)
				mensagens_enviadas_lock.Lock()
				mensagens_enviadas[0] += 1
				mensagens_enviadas_lock.Unlock()
				time.Sleep(time.Duration(4) * time.Second)
								
				if election_status_receive == 0 {
					send_message("2" + "|" + fmt.Sprint(process_id) + "|" + fmt.Sprint(process_rank) + ";", multicastPort)
					mensagens_enviadas_lock.Lock()
					mensagens_enviadas[2] += 1
					mensagens_enviadas_lock.Unlock()
				}				
			} else {
				lider_status_lock.Unlock()		
			}
	  	}
	}()

	// Imprime o lider atual a cada 10 segundos
	go func() {
		for true {
			time.Sleep(time.Duration(10) * time.Second)
			fmt.Println("")
			fmt.Println(lider_atual)
			fmt.Println("")
		}
	}()

	wg.Wait()
}