package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Command-line flags
	// default host: localhost, port: 6666, timeout: 30s
	host := flag.String("host", "localhost", "Server hostname")
	port := flag.Int("port", 6666, "Server port")
	timeout := flag.Int("timeout", 30, "Connection timeout in seconds")
	flag.Parse()

	// Connect to the server
	dialer := net.Dialer{Timeout: time.Duration(*timeout) * time.Second}
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to chat server at %s:%d\n", *host, *port)

	// Create a channel to listen for interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to listen for disconnect signals
	disconnectCh := make(chan struct{})
	
	// WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Receive messages from the server
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					// Normal disconnection
					break
				}
				fmt.Printf("\nConnection lost: %v\n", err)
				// Notify main thread that connection is closed
				close(disconnectCh)
				break
			}
			
			// Print the message from the server
				fmt.Print(message)
		}
	}()

	// Read user input and send messages to the server
	go func() {
		defer wg.Done()
		printHelp()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			


			
			// 发送命令到服务器 (Send command to server)
			_, err := conn.Write([]byte(line + "\n"))
			if err != nil {
				fmt.Printf("Error sending message: %v\n", err)
				close(disconnectCh)
				break
			}
			
			// 记录用户操作 (Log user actions)
			if strings.HasPrefix(line, "/NICK ") || strings.HasPrefix(line, "/N ") {
				nickname := strings.Split(line, " ")[1]
				fmt.Printf("Attempting to register/change nickname to: %s\n", nickname)
			}
		}

		if scanner.Err() == nil {
			fmt.Println("Disconnecting from server...")
			close(disconnectCh)
			return
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading from stdin: %v\n", err)
		}
		close(disconnectCh)
	}()

	// Wait for interrupt signal or disconnect signal
	select {
	case <-sigCh:
		fmt.Println("\nReceived interrupt signal")
	case <-disconnectCh:
		// Do nothing
	}

	// Close the connection
	conn.Close()

	// Wait for goroutines to finish
	wg.Wait()
	fmt.Println("Connection closed")
}

// printHelp prints the available commands to the user
func printHelp() {
	fmt.Println("Available Commands:")
	fmt.Println("  /NICK <nickname>, /N <nickname> - Set or change your nickname")
	fmt.Println("  /LIST, /L                       - List all connected users")
	fmt.Println("  /MSG <user> <message>, /M <user> <message> - Send a private message")
	fmt.Println("  /MSG * <message>, /M * <message>           - Send a message to all users")
	fmt.Println("-------------------")
}