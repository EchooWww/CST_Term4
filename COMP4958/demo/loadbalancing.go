package main

import (
	"bug-free/demo/mqueue" // Replace with your actual package path
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	// Demo 2: Load Balancing with Round-Robin

	fmt.Println("=== Load Balancing Demo ===")

	// Create data directory
	dataDir := "./data_lb_demo"
	os.MkdirAll(dataDir, 0755)
	defer os.RemoveAll(dataDir) // Clean up after demo

	// Create multiple nodes to demonstrate load balancing
	fmt.Println("1. Creating multiple nodes...")
	node1, err := mqueue.NewNode("node1", ":9001", dataDir)
	if err != nil {
		fmt.Printf("Error creating node1: %v\n", err)
		return
	}

	node2, err := mqueue.NewNode("node2", ":9002", dataDir)
	if err != nil {
		fmt.Printf("Error creating node2: %v\n", err)
		return
	}

	node3, err := mqueue.NewNode("node3", ":9003", dataDir)
	if err != nil {
		fmt.Printf("Error creating node3: %v\n", err)
		return
	}

	// Setup load balancer
	fmt.Println("3. Creating load balancer with round-robin strategy...")
	lb := mqueue.NewLoadBalancer()
	
	// IMPORTANT FIX: Use direct node connections instead of proxy connections
	node1Client, err := mqueue.DirectNodeClient("localhost:9001")
	if err != nil {
		fmt.Printf("Error connecting to node1: %v\n", err)
		return
	}
	
	node2Client, err := mqueue.DirectNodeClient("localhost:9002") 
	if err != nil {
		fmt.Printf("Error connecting to node2: %v\n", err)
		return
	}
	
	node3Client, err := mqueue.DirectNodeClient("localhost:9003")
	if err != nil {
		fmt.Printf("Error connecting to node3: %v\n", err)
		return
	}
	
	lb.AddNode("node1", node1Client)
	lb.AddNode("node2", node2Client)
	lb.AddNode("node3", node3Client)

	// Create the same queue on each node
	queueName := "distributed_queue"
	fmt.Println("4. Creating queue on all nodes:", queueName)
	node1.Broker.QueueManager.CreateQueue(queueName)
	node2.Broker.QueueManager.CreateQueue(queueName)
	node3.Broker.QueueManager.CreateQueue(queueName)

	// Setup message counters per node
	var counters sync.Map
	counters.Store("node1", 0)
	counters.Store("node2", 0)
	counters.Store("node3", 0)

	// Setup consumers on all nodes to track message distribution
	fmt.Println("5. Setting up consumers on all nodes...")
	
	consumer1 := mqueue.NewClient(node1.Broker)
	consumer1.Consume(queueName, func(msg mqueue.Message) {
		val, _ := counters.Load("node1")
		counters.Store("node1", val.(int)+1)
		fmt.Printf("   Node1 received: %s\n", string(msg.Body))
	})
	
	consumer2 := mqueue.NewClient(node2.Broker)
	consumer2.Consume(queueName, func(msg mqueue.Message) {
		val, _ := counters.Load("node2")
		counters.Store("node2", val.(int)+1)
		fmt.Printf("   Node2 received: %s\n", string(msg.Body))
	})
	
	consumer3 := mqueue.NewClient(node3.Broker)
	consumer3.Consume(queueName, func(msg mqueue.Message) {
		val, _ := counters.Load("node3")
		counters.Store("node3", val.(int)+1)
		fmt.Printf("   Node3 received: %s\n", string(msg.Body))
	})

	// Make sure consumers are ready before sending messages
	time.Sleep(time.Second)

	// Send messages through the load balancer
	fmt.Println("6. Sending messages through load balancer...")
	for i := 1; i <= 9; i++ {
		msg := fmt.Sprintf("Load balanced message #%d", i)
		fmt.Printf("   Sending: %s\n", msg)
		
		// Use load balancer to distribute messages
		err := lb.SendMessage(queueName, []byte(msg))
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		}
		
		// Wait to ensure message is processed
		time.Sleep(time.Millisecond * 300)
	}

	// Wait for message processing to complete
	time.Sleep(time.Second * 2)

	// Display message distribution
	fmt.Println("7. Message distribution across nodes:")
	
	node1Count, _ := counters.Load("node1")
	node2Count, _ := counters.Load("node2")
	node3Count, _ := counters.Load("node3")
	
	fmt.Printf("   Node1: %d messages\n", node1Count)
	fmt.Printf("   Node2: %d messages\n", node2Count)
	fmt.Printf("   Node3: %d messages\n", node3Count)
	
	totalMessages := node1Count.(int) + node2Count.(int) + node3Count.(int)
	fmt.Printf("8. Total messages processed: %d\n", totalMessages)
	
	fmt.Println("=== Load Balancing Demo Completed ===")
}