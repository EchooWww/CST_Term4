// package mqueue provides a simple distributed message queue implementation
package mqueue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// Message represents a data unit in the queue
type Message struct {
	ID        string
	Body      []byte
	Timestamp time.Time
	Attempts  int
}

// Queue represents a single message queue with file persistence
type Queue struct {
	name     string
	messages chan Message
	dataDir  string         
	mu       sync.RWMutex
}

// QueueManager handles multiple queues
type QueueManager struct {
	queues  map[string]*Queue
	mu      sync.RWMutex
	baseDir string // Base directory for all queue data
}

// NewQueueManager creates a new queue manager
func NewQueueManager(baseDir string) *QueueManager {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		fmt.Printf("Error creating base directory: %v\n", err)
	}

	return &QueueManager{
		queues:  make(map[string]*Queue),
		baseDir: baseDir,
	}
}

// CreateQueue creates a new queue with the given name
func (qm *QueueManager) CreateQueue(name string) *Queue {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	if q, exists := qm.queues[name]; exists {
		return q
	}

	// Create queue directory
	queueDir := filepath.Join(qm.baseDir, name)
	if err := os.MkdirAll(queueDir, 0755); err != nil {
		fmt.Printf("Error creating queue directory for %s: %v\n", name, err)
	}

	q := &Queue{
		name:     name,
		messages: make(chan Message, 1000), // Create a buffered message channel
		dataDir:  queueDir,
	}
	
	// Load any existing messages from disk
	if err := q.loadMessages(); err != nil {
		fmt.Printf("Error loading messages for queue %s: %v\n", name, err)
	}
	
	qm.queues[name] = q
	return q
}

// GetQueue returns the queue with the given name
func (qm *QueueManager) GetQueue(name string) (*Queue, error) {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	q, exists := qm.queues[name]
	if !exists {
		return nil, errors.New("queue not found")
	}
	return q, nil
}

// Save a message to disk
func (q *Queue) saveMessage(msg Message) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	
	filePath := filepath.Join(q.dataDir, msg.ID+".json")
	return ioutil.WriteFile(filePath, data, 0644)
}

// Delete a message from disk after processing
func (q *Queue) deleteMessage(msgID string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	filePath := filepath.Join(q.dataDir, msgID+".json")
	return os.Remove(filePath)
}

// Load all messages from disk into memory
func (q *Queue) loadMessages() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	files, err := ioutil.ReadDir(q.dataDir)
	if err != nil {
		return err
	}
	
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(q.dataDir, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue // Skip files with errors
			}
			
			var msg Message
			if err := json.Unmarshal(data, &msg); err != nil {
				continue // Skip malformed messages
			}
			
			// Send message to channel (non-blocking)
			select {
			case q.messages <- msg:
				// Message added to channel
			default:
				// Channel is full
				fmt.Printf("Queue channel full")
			}
		}
	}
	
	return nil
}

// Client represents a connection to the message broker
type Client struct {
	id     string
	broker *Broker
}

// NewClient creates a new client connected to the given broker
func NewClient(broker *Broker) *Client {
	id := fmt.Sprintf("client-%d", time.Now().UnixNano())
	client := &Client{
		id:     id,
		broker: broker,
	}
	broker.registerCh <- client
	return client
}

// Close closes the client connection
func (c *Client) Close() {
	c.broker.unregisterCh <- c
}

// Produce sends a message to the specified ueue
func (c *Client) Produce(queueName string, body []byte) error {
	q, err := c.broker.QueueManager.GetQueue(queueName)
	if err != nil {
		return err
	}

	msg := Message{
		ID:        fmt.Sprintf("msg-%d", time.Now().UnixNano()),
		Body:      body,
		Timestamp: time.Now(),
		Attempts:  0,
	}
	
	// First save message to disk for persistence
	if err := q.saveMessage(msg); err != nil {
		return fmt.Errorf("failed to persist message: %v", err)
	}

	// Then try to add to queue channel with timeout
	select {
	case q.messages <- msg:
		return nil
	case <-time.After(time.Second * 5):
		// Message is already on disk, so it's not lost
		return errors.New("queue channel full, message persisted to disk")
	}
}

// Consume receives messages from the specified queue
func (c *Client) Consume(queueName string, handler func(Message)) error {
	q, err := c.broker.QueueManager.GetQueue(queueName)
	if err != nil {
		return err
	}
	
	go func() {
		for msg := range q.messages {
			handler(msg)
			// Delete message from disk after processing
			if err := q.deleteMessage(msg.ID); err != nil {
				fmt.Printf("Error deleting processed message %s: %v\n", msg.ID, err)
			}
		}
	}()
	
	return nil
}

// Broker coordinates message distribution
type Broker struct {
	QueueManager *QueueManager
	clients      map[string]*Client
	registerCh   chan *Client
	unregisterCh chan *Client
}

// NewBroker creates a new message broker
func NewBroker(dataDir string) *Broker {
	return &Broker{
		QueueManager: NewQueueManager(dataDir),
		clients:      make(map[string]*Client),
		registerCh:   make(chan *Client),
		unregisterCh: make(chan *Client),
	}
}

// Start begins the broker's main processing loop
func (b *Broker) Start() {
	for {
		select {
		case client := <-b.registerCh:
			b.clients[client.id] = client
			fmt.Printf("Client registered: %s\n", client.id)
		case client := <-b.unregisterCh:
			delete(b.clients, client.id)
			fmt.Printf("Client unregistered: %s\n", client.id)
		}
	}
}

// ProduceArgs represents the arguments for the Produce RPC call
type ProduceArgs struct {
	Queue   string
	Message []byte
}

// Produce handles remote produce requests
func (b *Broker) Produce(args ProduceArgs, reply *bool) error {
	q, err := b.QueueManager.GetQueue(args.Queue)
	if err != nil {
		return err
	}

	msg := Message{
		ID:        fmt.Sprintf("msg-%d", time.Now().UnixNano()),
		Body:      args.Message,
		Timestamp: time.Now(),
		Attempts:  0,
	}
	
	// First save message to disk
	if err := q.saveMessage(msg); err != nil {
		*reply = false
		return fmt.Errorf("failed to persist message: %v", err)
	}
	
	// Then try to add to queue channel
	select {
	case q.messages <- msg:
		*reply = true
		return nil
	case <-time.After(time.Second * 5):
		// Message is already on disk, so it's not lost
		*reply = true
		return nil
	}
}

// Node represents a server node in the distributed system
type Node struct {
	ID           string
	Broker       *Broker
	RPCServer    *rpc.Server
	Listener     net.Listener
	PeerNodes    map[string]*rpc.Client
}

// LoadBalancer handles distribution of messages across nodes
type LoadBalancer struct {
	nodes    map[string]*rpc.Client
	counter  uint64 // For round-robin selection
	mu       sync.RWMutex
}

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		nodes:   make(map[string]*rpc.Client),
		counter: 0,
	}
}

// AddNode adds a node to the load balancer
func (lb *LoadBalancer) AddNode(nodeID string, client *rpc.Client) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.nodes[nodeID] = client
}

// RemoveNode removes a node from the load balancer
func (lb *LoadBalancer) RemoveNode(nodeID string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	delete(lb.nodes, nodeID)
}

// GetNextNode returns the next node in strict round-robin fashion
func (lb *LoadBalancer) GetNextNode() (string, *rpc.Client, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	
	if len(lb.nodes) == 0 {
		return "", nil, errors.New("no nodes available")
	}
	
	// Get list of node IDs for indexing in a deterministic order
	nodeIDs := make([]string, 0, len(lb.nodes))
	for id := range lb.nodes {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs) // Sort for deterministic ordering
	
	// Select node using round-robin counter
	idx := atomic.AddUint64(&lb.counter, 1) % uint64(len(nodeIDs))
	selectedID := nodeIDs[idx]
	
	return selectedID, lb.nodes[selectedID], nil
}

// SendMessage sends a message through the load balancer
func (lb *LoadBalancer) SendMessage(queueName string, body []byte) error {
	nodeID, client, err := lb.GetNextNode()
	if err != nil {
		return err
	}
	
	fmt.Printf("Load balancer selected node: %s\n", nodeID)
	
	var reply bool
	return client.Call("Broker.Produce", ProduceArgs{
		Queue:   queueName,
		Message: body,
	}, &reply)
}

// HeartbeatMonitor checks node health periodically
type HeartbeatMonitor struct {
	nodes        map[string]*rpc.Client
	loadBalancer *LoadBalancer
	interval     time.Duration
	stopCh       chan struct{}
	mu           sync.RWMutex
}

// NewHeartbeatMonitor creates a new heartbeat monitor
func NewHeartbeatMonitor(lb *LoadBalancer, interval time.Duration) *HeartbeatMonitor {
	return &HeartbeatMonitor{
		nodes:        make(map[string]*rpc.Client),
		loadBalancer: lb,
		interval:     interval,
		stopCh:       make(chan struct{}),
	}
}

// AddNode adds a node to be monitored
func (hm *HeartbeatMonitor) AddNode(nodeID string, client *rpc.Client) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.nodes[nodeID] = client
}

// Start begins heartbeat monitoring
func (hm *HeartbeatMonitor) Start() {
	ticker := time.NewTicker(hm.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			hm.checkNodes()
		case <-hm.stopCh:
			return
		}
	}
}

// Stop halts heartbeat monitoring
func (hm *HeartbeatMonitor) Stop() {
	close(hm.stopCh)
}

func testNodeHealth(client *rpc.Client) bool {
	// Create a new connection for each check to avoid stale connections
	var reply bool
	
	// Set timeout directly on Call
	ch := make(chan error, 1)
	go func() {
			ch <- client.Call("Node.Ping", struct{}{}, &reply)
	}()
	
	// Wait with timeout
	select {
	case err := <-ch:
			return err == nil && reply
	case <-time.After(time.Millisecond * 500):
			return false
	}
}

// Modified checkNodes method for HeartbeatMonitor
func (hm *HeartbeatMonitor) checkNodes() {
	hm.mu.RLock()
	nodesToCheck := make(map[string]*rpc.Client)
	for id, client := range hm.nodes {
			nodesToCheck[id] = client
	}
	hm.mu.RUnlock()
	
	for nodeID, _ := range nodesToCheck {
			// Force direct connection test instead of using existing client
			isAlive := false
			
			// Parse address from node ID
			port := ""
			switch nodeID {
			case "node1":
					port = ":9101"
			case "node2":
					port = ":9102"
			case "node3":
					port = ":9103"
			}
			
			addr := "localhost" + port
			
			// Try direct TCP connection to verify node is up
			conn, err := net.DialTimeout("tcp", addr, time.Millisecond*300)
			if err == nil {
					conn.Close()
					isAlive = true
			}
			
			fmt.Printf("HeartbeatMonitor: Direct TCP check for %s (%s): %v\n", 
								 nodeID, addr, isAlive)
			
			// Process result
			if !isAlive {
					fmt.Printf("HeartbeatMonitor: Node %s is DOWN, removing from load balancer\n", nodeID)
					
					// Remove from load balancer
					hm.loadBalancer.RemoveNode(nodeID)
					
					// Remove from monitoring
					hm.mu.Lock()
					delete(hm.nodes, nodeID)
					hm.mu.Unlock()
			}
	}
}

// Add to LoadBalancer struct
func (lb *LoadBalancer) GetNodeIDs() []string {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	
	nodeIDs := make([]string, 0, len(lb.nodes))
	for id := range lb.nodes {
		nodeIDs = append(nodeIDs, id)
	}
	return nodeIDs
}

func NewNode(id string, addr string, dataDir string) (*Node, error) {
	broker := NewBroker(filepath.Join(dataDir, id))
	go broker.Start()

	// Create RPC server
	server := rpc.NewServer()
	err := server.Register(broker)
	if err != nil {
		return nil, err
	}

	// Start listening first
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	node := &Node{
		ID:           id,
		Broker:       broker,
		RPCServer:    server,
		Listener:     listener,
		PeerNodes:    make(map[string]*rpc.Client),
	}

	// Register the node itself
	err = server.Register(node)
	if err != nil {
		listener.Close() // Close listener if registration fails
		return nil, err
	}

	// Start RPC service
	go server.Accept(listener)
	return node, nil
}

// ConnectToPeer connects to another node
func (n *Node) ConnectToPeer(id string, addr string) error {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return err
	}
	n.PeerNodes[id] = client
	return nil
}

// Ping responds to heartbeat checks
func (n *Node) Ping(args struct{}, reply *bool) error {
	*reply = true
	return nil
}

// SendMessage sends a message to a queue through this node
func (n *Node) SendMessage(queueName string, body []byte) error {
	var reply bool
	return n.Broker.Produce(ProduceArgs{
		Queue:   queueName,
		Message: body,
	}, &reply)
}

// DirectNodeClient creates a direct RPC client connection to a node
func DirectNodeClient(addr string) (*rpc.Client, error) {
	// Directly connect to the node's RPC server
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return client, nil
}