package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	privateIP := "0.0.0.0"
	port := "8090"
	listener, err := net.Listen("tcp", privateIP+":"+port)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port " + port)

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar a conexão:", err)
			continue
		}

		// Handle client connection in a goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	// Create a buffer to read data into
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Conexão fechada pelo cliente(EOF): %s\n", conn.RemoteAddr())
			} else {
				fmt.Printf("Erro ao ler os dados do cliente: %s\n", err)
			}
			return
		}

		//Novo tratamento
		data := buffer[:n]
		hexData := hex.EncodeToString(data)

		if len(hexData) == 68 { // 34 bytes = 68 characters in hex
			imeiAck := []byte{0x01}
			_, err := conn.Write(imeiAck)
			if err != nil {
				fmt.Println("Erro ao enviar o ACK do IMEI:", err)
				return
			}
			fmt.Println("Recebido o IMEI e retornado o ACK:", err)
		} else {
			// Parse e manipular dados
			fmt.Println("Brincar com os dados")

			// Contar o número de dados AVL
			numData := len(strings.Split(hexData, " ")) / 2
			hexString := fmt.Sprintf("%08X", numData)
			packetAck, err := hex.DecodeString(hexString)
			if err != nil {
				fmt.Println("Erro ao converter número de dados para hex:", err)
				return
			}
			fmt.Println("Dados após o ACK", packetAck)
			_, err = conn.Write(packetAck)
			if err != nil {
				fmt.Println("Erro ao enviar o ACK do pacote:", err)
				return
			}
		}

		// if n > 0 {
		// 	message := string(buffer[:n])
		// 	fmt.Printf("Mensagem recebida do cliente %s: %s\n", conn.RemoteAddr(), message)

		// 	// Aqui você pode processar a mensagem conforme necessário
		// 	// Por exemplo, você pode analisar os dados da mensagem para extrair as informações do rastreador GPS
		// }
	}
}
