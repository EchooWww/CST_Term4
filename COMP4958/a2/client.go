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
	// 解析命令行参数
	host := flag.String("host", "localhost", "Server hostname")
	port := flag.Int("port", 6666, "Server port")
	timeout := flag.Int("timeout", 30, "Connection timeout in seconds")
	flag.Parse()

	// 连接到服务器，带超时
	dialer := net.Dialer{Timeout: time.Duration(*timeout) * time.Second}
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to chat server at %s:%d\n", *host, *port)
	fmt.Println("Type /help for available commands")

	// 设置信号处理以优雅地关闭连接
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 创建一个通道来同步连接终止
	disconnectCh := make(chan struct{})
	
	// 创建一个等待组来等待 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(2)

	// 接收来自服务器的消息
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					// 正常关闭 (Normal closure)
					break
				}
				fmt.Printf("\nConnection lost: %v\n", err)
				// 通知主线程连接已关闭 (Notify main thread that connection is closed)
				close(disconnectCh)
				break
			}
			
			// 处理特定响应消息 (Handle specific response messages)
			trimmedMsg := strings.TrimSpace(message)
			if strings.Contains(trimmedMsg, "Nickname") && 
			   (strings.Contains(trimmedMsg, "registered") || 
			    strings.Contains(trimmedMsg, "changed") || 
			    strings.Contains(trimmedMsg, "already in use")) {
				fmt.Printf("\n[STATUS] %s\n", trimmedMsg)
			} else {
				// 打印普通消息 (Print regular messages)
				fmt.Print(message)
			}
		}
	}()

	// 发送用户输入到服务器
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			
			// 客户端处理特殊命令
			if line == "/help" {
				printHelp()
				continue
			}
			if line == "/quit" || line == "/exit" {
				fmt.Println("Disconnecting from server...")
				close(disconnectCh)
				return
			}
			
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

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading from stdin: %v\n", err)
		}
		close(disconnectCh)
	}()

	// 等待中断信号或连接断开
	select {
	case <-sigCh:
		fmt.Println("\nReceived interrupt signal")
	case <-disconnectCh:
		// 其他 goroutine 已经通知断开连接
	}

	// 关闭连接
	conn.Close()

	// 等待所有 goroutine 完成
	wg.Wait()
	fmt.Println("Connection closed")
}

// 打印帮助信息
func printHelp() {
	fmt.Println("\n--- Chat Client Help ---")
	fmt.Println("Server Commands:")
	fmt.Println("  /NICK <nickname>, /N <nickname> - Set or change your nickname")
	fmt.Println("  /LIST, /L                       - List all connected users")
	fmt.Println("  /MSG <user> <message>, /M <user> <message> - Send a private message")
	fmt.Println("  /MSG * <message>, /M * <message>           - Send a message to all users")
	fmt.Println("\nClient Commands:")
	fmt.Println("  /help  - Show this help message")
	fmt.Println("  /quit  - Disconnect from the server")
	fmt.Println("  /exit  - Disconnect from the server")
	fmt.Println("-------------------")
}