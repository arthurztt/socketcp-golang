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

	// Goroutine separada para receber as mensagens enquanto o usuário digita
	go func(){
		scanner := bufio.NewScanner(conn)
		for scanner.Scan(){
			fmt.Printf("\r%s\n> ", scanner.Text())
			scanner.Err()
		}
		fmt.Println("\nServidor encerrou a conexão.")
		os.Exit(0)
	}()

	// Loop para enviar as mensagens
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

}