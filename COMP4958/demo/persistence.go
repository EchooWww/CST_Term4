package main

import (
	"bug-free/demo/mqueue" // Replace with your actual package path
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Demo 1: Message Persistence

	fmt.Println("=== Message Persistence Demo ===")

	// Create data directory
	dataDir := "./data_persistence_demo"
	os.MkdirAll(dataDir, 0755)
	
	// Create broker with persistence
	fmt.Println("1. Creating broker with persistence...")
	broker := mqueue.NewBroker(dataDir)
	go broker.Start()

	// Create queue
	queueName := "persistence_test"
	fmt.Println("2. Creating queue:", queueName)
	broker.QueueManager.CreateQueue(queueName)

	// Create producer
	producer := mqueue.NewClient(broker)
	defer producer.Close()

	// Send messages to queue - persisted to disk
	fmt.Println("3. Sending messages that will be persisted to disk...")
	for i := 1; i <= 5; i++ {
		msg := fmt.Sprintf("Persistent message #%d", i)
		fmt.Printf("   Producing: %s\n", msg)
		producer.Produce(queueName, []byte(msg))
		time.Sleep(time.Millisecond * 100)
	}

	// Check files on disk
	queueDir := fmt.Sprintf("%s/%s", dataDir, queueName)
	files, _ := os.ReadDir(queueDir)
	fmt.Printf("4. Disk persistence check: %d files saved in %s\n", len(files), queueDir)
	
	// List files for demonstration
	fmt.Println("   Files on disk:")
	for _, file := range files {
		// Read sample of file content
		filePath := filepath.Join(queueDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err == nil && len(content) > 0 {
			// Show first 50 chars of content
			contentPreview := string(content)
			if len(contentPreview) > 50 {
				contentPreview = contentPreview[:50] + "..."
			}
			fmt.Printf("   - %s: %s\n", file.Name(), contentPreview)
		} else {
			fmt.Printf("   - %s\n", file.Name())
		}
	}
	
	// Pause for file inspection
	fmt.Println("5. Pausing for 3 seconds to allow manual inspection of files...")
	fmt.Printf("   Check files at: %s\n", queueDir)
	time.Sleep(time.Second * 3)

	// Simulate broker restart
	fmt.Println("6. Simulating broker restart...")
	
	// Create new broker that loads messages from disk
	newBroker := mqueue.NewBroker(dataDir)
	go newBroker.Start()

	// Recreate queue - should load existing messages
	fmt.Println("7. Recreating queue after restart...")
	newBroker.QueueManager.CreateQueue(queueName)

	// Setup consumer to read messages
	fmt.Println("8. Setting up consumer to read recovered messages...")
	consumer := mqueue.NewClient(newBroker)
	defer consumer.Close()

	messageCount := 0
	receivedMsgs := make(chan string, 10)

	consumer.Consume(queueName, func(msg mqueue.Message) {
		messageCount++
		receivedMsgs <- string(msg.Body)
	})

	// Wait to receive all messages
	time.Sleep(time.Second)
	close(receivedMsgs)

	// Display recovered messages
	fmt.Println("9. Messages recovered after restart:")
	for msg := range receivedMsgs {
		fmt.Printf("   Recovered: %s\n", msg)
	}

	fmt.Printf("10. Recovery complete: %d messages restored from disk\n", messageCount)
	
	// Count remaining files
	remainingFiles, _ := os.ReadDir(queueDir)
	fmt.Printf("11. Remaining files after consumption: %d\n", len(remainingFiles))
	
	// Final message with delay
	fmt.Println("12. Demo completed - files remain at:", queueDir)
	fmt.Println("    Program will exit in 3 seconds...")
	
	// Use a timeout instead of blocking forever
	time.Sleep(time.Second * 3)
	fmt.Println("    Demo exit.")
}