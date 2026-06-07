package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var (
	clients = make(map[net.Conn]bool)
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

		// Registra um novo cliente
		clientsMu.Lock()
		clients[conn] = true
		clientsMu.Unlock()
		
		fmt.Printf("Cliente conectado: %s | Total: %d\n", conn.RemoteAddr(), len(clients))

		// Trata cada cliente em uma goroutine separada
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		// Remove clientes ao desconectar
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()

		fmt.Printf("Cliente desconectado: %s | Total %d\n", conn.RemoteAddr(), len(clients))
		conn.Close()
	}()

	
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Printf("[%s]: %s\n", conn.RemoteAddr(), msg)
		broadcast(msg, conn)

		
		if scanner.Err() != nil {
			fmt.Println("Erro ao ler mensagem:", scanner.Err())
			return
		}
	}
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