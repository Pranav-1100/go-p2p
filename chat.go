package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func isExitCommand(msg string) bool {
	cmd := strings.ToLower(strings.TrimSpace(msg))
	return cmd == "exit" || cmd == "quit"
}

func handleChat(conn net.Conn) {
	defer conn.Close()
	done := make(chan bool)

	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Println("\n✗ Peer disconnected")
				}
				done <- true
				return
			}
			message = strings.TrimSpace(message)
			if message != "" {
				fmt.Printf("Peer: %s\n", message)
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			message := scanner.Text()

			if isExitCommand(message) {
				fmt.Println("\n✓ Closing connection...")
				done <- true
				return
			}

			_, err := fmt.Fprintf(conn, "%s\n", message)
			if err != nil {
				fmt.Printf("\n✗ Failed to send: %v\n", err)
				done <- true
				return
			}
			fmt.Printf("You: %s\n", message)
		}
		if err := scanner.Err(); err != nil {
			done <- true
		}
	}()

	<-done
	fmt.Println("Chat ended.")
}