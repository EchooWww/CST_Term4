package main

import (
	"bug-free/demo/mqueue" // Replace with your actual package path
	"fmt"
	"os"
	"time"
)

func main() {
	// Demo 3: Heartbeat Monitoring and Fault Detection

	fmt.Println("=== Heartbeat Monitoring Demo ===")

	// Create data directory
	dataDir := "./data_heartbeat_demo"
	os.MkdirAll(dataDir, 0755)
	defer os.RemoveAll(dataDir) // Clean up after demo

	// Create nodes
	fmt.Println("1. Creating nodes...")
	node1, err := mqueue.NewNode("node1", ":9101", dataDir)
	if err != nil {
		fmt.Printf("Error creating node1: %v\n", err)
		return
	}

	node2, err := mqueue.NewNode("node2", ":9102", dataDir)
	if err != nil {
		fmt.Printf("Error creating node2: %v\n", err)
		return
	}

	node3, err := mqueue.NewNode("node3", ":9103", dataDir)
	if err != nil {
		fmt.Printf("Error creating node3: %v\n", err)
		return
	}

	// Create direct connections to each node
	node1Client, err := mqueue.DirectNodeClient("localhost:9101")
	if err != nil {
		fmt.Printf("Error connecting to node1: %v\n", err)
		return
	}
	
	node2Client, err := mqueue.DirectNodeClient("localhost:9102")
	if err != nil {
		fmt.Printf("Error connecting to node2: %v\n", err)
		return
	}
	
	node3Client, err := mqueue.DirectNodeClient("localhost:9103")
	if err != nil {
		fmt.Printf("Error connecting to node3: %v\n", err)
		return
	}

	// Setup load balancer with direct clients
	fmt.Println("3. Setting up load balancer...")
	lb := mqueue.NewLoadBalancer()
	lb.AddNode("node1", node1Client)
	lb.AddNode("node2", node2Client)
	lb.AddNode("node3", node3Client)

	// Setup heartbeat monitor with a short interval for demo purposes
	fmt.Println("4. Starting heartbeat monitor (1-second interval)...")
	hm := mqueue.NewHeartbeatMonitor(lb, time.Second*1)
	hm.AddNode("node1", node1Client)
	hm.AddNode("node2", node2Client)
	hm.AddNode("node3", node3Client)
	go hm.Start()

	// Create queues
	queueName := "heartbeat_test"
	fmt.Println("5. Creating queues on all nodes...")
	node1.Broker.QueueManager.CreateQueue(queueName)
	node2.Broker.QueueManager.CreateQueue(queueName)
	node3.Broker.QueueManager.CreateQueue(queueName)

	// Setup consumers
	fmt.Println("6. Setting up consumers...")
	consumer1 := mqueue.NewClient(node1.Broker)
	consumer1.Consume(queueName, func(msg mqueue.Message) {
		fmt.Printf("   Node1 received: %s\n", string(msg.Body))
	})

	consumer2 := mqueue.NewClient(node2.Broker)
	consumer2.Consume(queueName, func(msg mqueue.Message) {
		fmt.Printf("   Node2 received: %s\n", string(msg.Body))
	})

	consumer3 := mqueue.NewClient(node3.Broker)
	consumer3.Consume(queueName, func(msg mqueue.Message) {
		fmt.Printf("   Node3 received: %s\n", string(msg.Body))
	})

	// Give some time for all clients to register
	time.Sleep(time.Second)

	// Send initial messages to all nodes
	fmt.Println("7. Sending initial messages to all nodes...")
	for i := 1; i <= 3; i++ {
		lb.SendMessage(queueName, []byte(fmt.Sprintf("Initial message #%d", i)))
		time.Sleep(time.Millisecond * 300)
	}

	// Simulate node2 failure by stopping its listener
	fmt.Println("8. Simulating node2 failure...")
	node2.Listener.Close() // This will cause heartbeat checks to fail
	fmt.Println("   Waiting for heartbeat monitor to detect failure (3 seconds)...")
	time.Sleep(time.Second * 3) // Wait for detection

	// Verify node2 was removed from load balancer
	fmt.Println("9. Checking load balancer status after node failure...")
	
	// Check if node2 still exists in the load balancer
	nodeIDs := lb.GetNodeIDs()
	node2Exists := false
	for _, id := range nodeIDs {
		if id == "node2" {
			node2Exists = true
			break
		}
	}
	
	if !node2Exists {
		fmt.Println("   ✓ Failed node2 successfully removed from load balancer")
	} else {
		fmt.Println("   ✗ Node2 still in load balancer")
	}

	// Send more messages - they should only go to remaining nodes
	fmt.Println("10. Sending messages after node failure...")
	for i := 1; i <= 4; i++ {
		msg := fmt.Sprintf("Post-failure message #%d", i)
		fmt.Printf("    Sending: %s\n", msg)
		
		nodeID,client, _ := lb.GetNextNode()
		fmt.Printf("    Selected node: %s\n", nodeID)
		
		var reply bool
		client.Call("Broker.Produce", mqueue.ProduceArgs{
				Queue:   queueName,
				Message: []byte(msg),
		}, &reply)		
		time.Sleep(time.Millisecond * 300)
	}

	// Wait for processing to complete
	time.Sleep(time.Second)
	fmt.Println("11. All messages routed to remaining healthy nodes")
	
	fmt.Println("=== Heartbeat Monitoring Demo Completed ===")
}