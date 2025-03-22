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

// Client 表示一个连接的客户端
type Client struct {
	conn     net.Conn    // TCP连接
	nickname string      // 客户端昵称
	outCh    chan string // 发送消息的通道
}

// Server 管理客户端和消息路由
type Server struct {
	mu      sync.RWMutex        // 保护 clients map 的互斥锁
	clients map[string]*Client  // 昵称到客户端的映射
	logger  *log.Logger         // 日志记录器
}

// 创建一个新的聊天服务器实例
func NewServer(logger *log.Logger) *Server {
	return &Server{
		clients: make(map[string]*Client),
		logger:  logger,
	}
}

// 注册新客户端
func (s *Server) RegisterClient(nickname string, client *Client) (bool, string) {
	// 验证昵称格式
	valid, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]{0,11}$`, nickname)
	if !valid {
		return false, "Invalid nickname format"
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[nickname]; exists {
		return false, fmt.Sprintf("Nickname %s already in use", nickname)
	}

	s.clients[nickname] = client
	s.logger.Printf("User registered with nickname: %s", nickname)
	return true, fmt.Sprintf("Nickname %s registered successfully", nickname)
}

// 注销客户端
func (s *Server) UnregisterClient(nickname string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[nickname]; exists {
		delete(s.clients, nickname)
		s.logger.Printf("User %s left the chat", nickname)
	}
}

// 更改客户端昵称
func (s *Server) ChangeNickname(oldNick, newNick string, client *Client) (bool, string) {
	// 验证新昵称
	valid, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]{0,11}$`, newNick)
	if !valid {
		return false, "Invalid nickname format"
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[newNick]; exists {
		return false, fmt.Sprintf("Nickname %s already in use", newNick)
	}

	// 检查旧昵称是否存在
	if oldClient, exists := s.clients[oldNick]; exists && oldClient == client {
		delete(s.clients, oldNick)
		s.clients[newNick] = client
		s.logger.Printf("User %s changed nickname to %s", oldNick, newNick)
		return true, fmt.Sprintf("Nickname changed to %s successfully", newNick)
	}

	return false, "Cannot change nickname: current nickname not found"
}

// 列出所有注册用户
func (s *Server) ListUsers() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]string, 0, len(s.clients))
	for nick := range s.clients {
		users = append(users, nick)
	}
	
	// 按字母顺序排序用户列表
	sort.Strings(users)
	return users
}

// 发送消息给指定收件人
func (s *Server) SendMessage(sender, recipients, message string) ([]string, []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var recipientList []string
	if recipients == "*" {
		// 发送给除自己外的所有用户
		for nick := range s.clients {
			if nick != sender {
				recipientList = append(recipientList, nick)
			}
		}
	} else {
		// 解析收件人列表
		recipientList = strings.Split(recipients, ",")
		for i, r := range recipientList {
			recipientList[i] = strings.TrimSpace(r)
		}
	}

	success := []string{}
	failed := []string{}

	// 发送消息给所有收件人
	for _, r := range recipientList {
		if client, exists := s.clients[r]; exists {
			formattedMsg := fmt.Sprintf("%s: %s\r\n", sender, message)
			select {
			case client.outCh <- formattedMsg:
				success = append(success, r)
			default:
				// 如果通道已满，记为失败
				s.logger.Printf("Failed to send message to %s: channel full", r)
				failed = append(failed, r)
			}
		} else {
			failed = append(failed, r)
		}
	}

	return success, failed
}

// 处理客户端连接
func handleConnection(server *Server, conn net.Conn) {
	defer func() {
		server.logger.Printf("Connection from %s closed", conn.RemoteAddr())
		conn.Close()
	}()

	server.logger.Printf("New connection from %s", conn.RemoteAddr())

	// 初始化客户端
	client := &Client{
		conn:  conn,
		outCh: make(chan string, 10), // 使用缓冲通道
	}

	// 启动一个 goroutine 处理发送消息到客户端
	go func() {
		for msg := range client.outCh {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				server.logger.Printf("Error writing to client: %v", err)
				return
			}
		}
	}()

	// 发送欢迎消息
	conn.Write([]byte("Welcome to the Go Chat Server!\r\n"))
	conn.Write([]byte("Please set a nickname with /NICK <nickname> or /N <nickname>\r\n"))

	scanner := bufio.NewScanner(conn)
	var nickname string

	for scanner.Scan() {
		command := strings.TrimSpace(scanner.Text())
		var response string
		

		// 处理各种命令
		if strings.HasPrefix(command, "/NICK ") || strings.HasPrefix(command, "/N ") {
			// NICK 命令
			parts := strings.SplitN(command, " ", 2)
			if len(parts) < 2 {
				response = "Invalid NICK format. Use /NICK <nickname>"
			} else {
				newNick := strings.Split(parts[1], " ")[0] // 获取第一个参数
				if nickname == "" {
					// 注册新昵称
					success, msg := server.RegisterClient(newNick, client)
					response = msg
					if success {
							nickname = newNick
							client.nickname = nickname
					}
				} else {
					// 更改昵称
					success, msg := server.ChangeNickname(nickname, newNick, client); 
					response = msg
					if success {
						nickname = newNick
						client.nickname = nickname
					}
				}
			}
		} else if command == "/NICK" || command == "/N" {
			// 无参数的NICK命令
			response = "Invalid NICK format. Use /NICK <nickname>"
		} else if command == "/LIST" || command == "/L" || strings.HasPrefix(command, "/LIST ") || strings.HasPrefix(command, "/L ") {
			// LIST 命令（忽略额外参数）
			users := server.ListUsers()
			if len(users) == 0 {
				response = "No users currently connected."
			} else {
				response = "Users: " + strings.Join(users, ", ")
			}
		} else if strings.HasPrefix(command, "/MSG ") || strings.HasPrefix(command, "/M ") {
			// MSG 命令
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
			// 无效的MSG命令
			response = "Invalid MSG format. Use /MSG <recipients> <message>"
		} else {
			response = "Unknown command. Available commands: /NICK <nickname>, /LIST, /MSG <recipients> <message>"
		}

		// 发送响应并记录 (Send response and log it)
		conn.Write([]byte(response + "\r\n"))
	}

	if err := scanner.Err(); err != nil {
		server.logger.Printf("Error reading from client: %v", err)
	}

	// 客户端断开连接，清理资源
	if nickname != "" {
		server.UnregisterClient(nickname)
	}
	close(client.outCh)
}

func main() {
	// 命令行参数
	port := flag.Int("port", 6666, "Port to listen on")
	logFile := flag.String("log", "", "Log file (default: stdout)")
	flag.Parse()

	// 设置日志
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

	// 启动服务器
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	logger.Printf("Chat server started on port %d", *port)

	// 接受连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(server, conn)
	}
}