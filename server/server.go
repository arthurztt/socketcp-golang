package main

import (
	"bufio"
	"fmt"
	"net"
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

		fmt.Printf("Cliente conectado: %s\n", conn.RemoteAddr())

		// Trata cada cliente em uma goroutine separada
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Printf("[%s] Mensagem recebida: %s\n", conn.RemoteAddr(), msg)
		if scanner.Err() != nil {
			fmt.Println("Erro ao ler mensagem:", scanner.Err())
			return
		}
	}

	fmt.Printf("Cliente desconectado: %s\n", conn.RemoteAddr())
}