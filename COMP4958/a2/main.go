package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// Client stands for a chat client
type Client struct {
	conn     net.Conn    // TCP connection between client and server
	nickname string      // User nickname
	outCh    chan string // Channel to send messages
}

// Server to manage clients and connections
type Server struct {
	mu      sync.RWMutex        // mutex to protect clients map
	clients map[string]*Client  // map of usernames to clients
	logger  *log.Logger         // logger for server
}

// NewServer creates a new server instance
func NewServer(logger *log.Logger) *Server {
	return &Server{
		clients: make(map[string]*Client),
		logger:  logger,
	}
}

// RegisterClient registers a new client with the server
func (s *Server) RegisterClient(nickname string, client *Client) (bool, string) {
	// Validate nickname format: ensure it starts with a letter and contains only letters, numbers, and underscores
	valid, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]{0,11}$`, nickname)
	if !valid {
		return false, "Invalid nickname format"
	}

	// Lock the clients map to prevent concurrent access   
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the nickname is already in use
	if _, exists := s.clients[nickname]; exists {
		return false, fmt.Sprintf("Nickname %s already in use", nickname)
	}

	s.clients[nickname] = client
	s.logger.Printf("User registered with nickname: %s", nickname)
	return true, fmt.Sprintf("Nickname %s registered successfully", nickname)
}

// UnregisterClient removes a client from the server
func (s *Server) UnregisterClient(nickname string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[nickname]; exists {
		delete(s.clients, nickname)
		s.logger.Printf("User %s left the chat", nickname)
	}
}

// ChangeNickname changes a client's nickname
func (s *Server) ChangeNickname(oldNick, newNick string, client *Client) (bool, string) {
	// Validate nickname format: ensure it starts with a letter and contains only letters, numbers, and underscores
	valid, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]{0,11}$`, newNick)
	if !valid {
		return false, "Invalid nickname format"
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the new nickname is already in use
	if _, exists := s.clients[newNick]; exists {
		return false, fmt.Sprintf("Nickname %s already in use", newNick)
	}

	// Check if the old nickname exists and matches the client
	if oldClient, exists := s.clients[oldNick]; exists && oldClient == client {
		delete(s.clients, oldNick)
		s.clients[newNick] = client
		s.logger.Printf("User %s changed nickname to %s", oldNick, newNick)
		return true, fmt.Sprintf("Nickname changed to %s successfully", newNick)
	}

	return false, "Cannot change nickname: current nickname not found"
}

// ListUsers returns a list of all connected users
func (s *Server) ListUsers() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]string, 0, len(s.clients))
	for nick := range s.clients {
		users = append(users, nick)
	}
	
	// Sort the list of users alphabetically
	sort.Strings(users)
	return users
}

// SendMessage sends a message from a sender to one or more recipients
func (s *Server) SendMessage(sender, recipients, message string) ([]string, []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var recipientList []string
	if recipients == "*" {
		// Send to all users except the sender
		for nick := range s.clients {
			if nick != sender {
				recipientList = append(recipientList, nick)
			}
		}
	} else {
		// Split the recipients string by commas
		recipientList = strings.Split(recipients, ",")
		for i, r := range recipientList {
			recipientList[i] = strings.TrimSpace(r)
		}
	}

	success := []string{}
	failed := []string{}

	// Send the message to each recipient
	for _, r := range recipientList {
		if client, exists := s.clients[r]; exists {
			formattedMsg := fmt.Sprintf("%s: %s\r\n", sender, message)
			select {
			case client.outCh <- formattedMsg:
				success = append(success, r)
			default:
				// If the client's outCh is full, the message is dropped
				s.logger.Printf("Failed to send message to %s: channel full", r)
				failed = append(failed, r)
			}
		} else {
			failed = append(failed, r)
		}
	}

	return success, failed
}

// handleConnection handles a new client connection
func handleConnection(server *Server, conn net.Conn) {
	defer func() {
		server.logger.Printf("Connection from %s closed", conn.RemoteAddr())
		conn.Close()
	}()

	server.logger.Printf("New connection from %s", conn.RemoteAddr())

	// Initialize a new client
	client := &Client{
		conn:  conn,
		outCh: make(chan string, 10), // Use a buffered channel of capacity 10
	}

	// Start a goroutine to send messages to the client
	go func() {
		for msg := range client.outCh {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				server.logger.Printf("Error writing to client: %v", err)
				return
			}
		}
	}()

	// Send a welcome message to the client
	conn.Write([]byte("Welcome to the Go Chat Server!\r\n"))
	conn.Write([]byte("Please set a nickname with /NICK <nickname> or /N <nickname>\r\n"))

	scanner := bufio.NewScanner(conn)
	var nickname string

	for scanner.Scan() {
		command := strings.TrimSpace(scanner.Text())
		var response string
		

		// Handle commands
		if strings.HasPrefix(command, "/NICK ") || strings.HasPrefix(command, "/N ") {
			// NICK Command
			parts := strings.SplitN(command, " ", 2)
			if len(parts) < 2 {
				response = "Invalid NICK format. Use /NICK <nickname>"
			} else {
				newNick := strings.Split(parts[1], " ")[0] // Extract the nickname
				if nickname == "" {
					// register new nickname
					success, msg := server.RegisterClient(newNick, client)
					response = msg
					if success {
							nickname = newNick
							client.nickname = nickname
					}
				} else {
					// change nickname
					success, msg := server.ChangeNickname(nickname, newNick, client); 
					response = msg
					if success {
						nickname = newNick
						client.nickname = nickname
					}
				}
			}
		} else if command == "/NICK" || command == "/N" {
			// Invalid NICK command
			response = "Invalid NICK format. Use /NICK <nickname>"
		} else if command == "/LIST" || command == "/L" || strings.HasPrefix(command, "/LIST ") || strings.HasPrefix(command, "/L ") {
			// LIST Command
			users := server.ListUsers()
			if len(users) == 0 {
				response = "No users currently connected."
			} else {
				response = "Users: " + strings.Join(users, ", ")
			}
		} else if strings.HasPrefix(command, "/MSG ") || strings.HasPrefix(command, "/M ") {
			// MSG Command
			if nickname == "" {
				response = "You must set a nickname before sending messages. Use /NICK <nickname>"
			} else {
				parts := strings.SplitN(command, " ", 3)
				if len(parts) < 3 {
					response = "Invalid MSG format. Use /MSG <recipients> <message>"
				} else {
					recipients := parts[1]
					message := parts[2]
					success, failed := server.SendMessage(nickname, recipients, message)

					if len(failed) == 0 {
						recipientDisplay := recipients
						if recipients == "*" {
							recipientDisplay = "all users"
						}
						response = fmt.Sprintf("Message sent to %s", recipientDisplay)
					} else if len(success) == 0 {
						response = fmt.Sprintf("No recipients found: %s", strings.Join(failed, ", "))
					} else {
						response = fmt.Sprintf("Message sent to %s, recipients not found: %s",
							strings.Join(success, ", "), strings.Join(failed, ", "))
					}
				}
			}
		} else if strings.HasPrefix(command, "/MSG") || strings.HasPrefix(command, "/M") {
			// Invalid MSG command
			response = "Invalid MSG format. Use /MSG <recipients> <message>"
		} else {
			response = "Unknown command. Available commands: /NICK <nickname>, /LIST, /MSG <recipients> <message>"
		}

		// Send the response to the client 
		conn.Write([]byte(response + "\r\n"))
	}

	if err := scanner.Err(); err != nil {
		server.logger.Printf("Error reading from client: %v", err)
	}

	// Unregister the client when the connection is closed
	if nickname != "" {
		server.UnregisterClient(nickname)
	}
	close(client.outCh)
}

func main() {
	// Command-line flags
	// default port is 6666
	port := flag.Int("port", 6666, "Port to listen on")
	// default log file is stdout
	logFile := flag.String("log", "", "Log file (default: stdout)")
	flag.Parse()

	// Create a logger
	var logger *log.Logger
	if *logFile != "" {
		file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer file.Close()
		logger = log.New(file, "", log.LstdFlags)
	} else {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	server := NewServer(logger)

	// Start the server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	logger.Printf("Chat server started on port %d", *port)

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(server, conn)
	}
}