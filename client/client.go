package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Conecta ao servidor na porta 3535
	conn, err := net.Dial("tcp", "localhost:3535")
	if err != nil {
		fmt.Println("Erro ao conectar:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Conectado ao servidor! Digite suas mensagens (CTRL+C para sair):")

	scanner := bufio.NewScanner(os.Stdin)
	

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		scanner.Err()
		
		msg := scanner.Text()

		// Envia a mensagem com quebra de linha para o servidor
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			fmt.Println("Erro ao enviar mensagem:", err)
			break
		}
	}

	fmt.Println("Conexão encerrada.")
}