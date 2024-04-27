package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

func messageHandler(w http.ResponseWriter, r *http.Request) {
	// Configura o cabeçalho da resposta para indicar o tipo de conteúdo
	w.Header().Set("Content-Type", "text/plain")

	// Lê o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusInternalServerError)
		return
	}

	// Converte o corpo para string e escreve na resposta
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Mensagem recebida: %s\n", string(body))
	fmt.Println("Mensagem:" + string(body))
}

func main() {
	privateIP := "0.0.0.0"
	// Configura o handler para a rota "/"
	port := "8090"
	// http.HandleFunc("/", messageHandler)

	// // Inicia o servidor na porta 8080
	// fmt.Println("Servidor iniciado na porta " + port)
	// if err := http.ListenAndServe(":"+port, nil); err != nil {
	// 	fmt.Printf("Erro ao iniciar o servidor: %s\n", err)
	// }

	// Listen for incoming connections
	// listener, err := net.Listen("tcp", "54.210.84.139:"+port)
	listener, err := net.Listen("tcp", privateIP+":"+port)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port " + port)

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
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
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if n > 0 {
			// Process and use the data (here, we'll just print it)
			_, err = fmt.Println("Mensagem recebida:", string(buffer[:n]))
			if err != nil {
				fmt.Println("Erro ao ler o buffer:", err)
				return
			}
		}
	}
}

func getPrivateIP() (string, error) {
	privateIP := os.Getenv("PRIVATE_IP")
	if privateIP == "" {
		return "", fmt.Errorf("variável de ambiente PRIVATE_IP não está definida")
	}
	return privateIP, nil
}
