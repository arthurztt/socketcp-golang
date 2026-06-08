package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	conn 	 net.Conn
	username string
}

var (
	clients = make(map[net.Conn]*Client)
	clientsMu sync.Mutex
)

func main() {
	// Cria o listener TCP na porta 3535
	listener, err := net.Listen("tcp", ":3535")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor aguardando conexões na porta 3535...")

	for {
		// Bloqueia até um cliente conectar
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}

		fmt.Printf("Cliente conectado: %s | Total: %d\n", conn.RemoteAddr(), len(clients))

		// Trata cada cliente em uma goroutine separada
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Err()
	
	// Pede o username como primeira interação
	fmt.Fprintln(conn, "Digite seu username:")
	if !scanner.Scan(){
		return
	}
	username :=	scanner.Text()

	//Registra client com username
	client := &Client{conn: conn, username: username}

	clientsMu.Lock()
	clients[conn] = client
	clientsMu.Unlock()

	// Log no servidor
	fmt.Printf("%s entrou no chat [%s]\n", username, conn.RemoteAddr())

	// Avisa aos outros que alguém entrou
	broadcast(fmt.Sprintf(">>> %s entrou no chat!", username), conn)

	// Confirma para o próprio cliente
	fmt.Fprintf(conn, "Bem-vindo, %s! Você está no chat.\n", username)

	// Loop de mensagens
	for scanner.Scan() {
		msg := scanner.Text()

		formatted := fmt.Sprintf("[%s]: %s", username, msg)
	
	 	fmt.Println(formatted)

		broadcast(formatted, conn)
	}

	// Ao desconectar o servidor
	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()

	fmt.Printf("%s saiu do chat\n", username)
	broadcast(fmt.Sprintf(">>> %s saiu do chat.", username), conn)
}
func broadcast(msg string, sender net.Conn){
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn :=	range clients {
		if conn != sender { // Não reenvia para quem mandou
			fmt.Fprintln(conn, msg)
		}
	}
}
