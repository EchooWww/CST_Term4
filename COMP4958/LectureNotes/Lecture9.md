# Lecture 9: Concurrency in Go

## 1. Sync.Cond: Condition Variables

Condition variables allow goroutines to wait for or signal a condition.

Key characteristics:

- Works with a `sync.Mutex` to provide goroutine synchronization
- Three main methods:
  - `Wait()`: Blocks until condition is signaled
  - `Signal()`: Notifies one waiting goroutine
  - `Broadcast()`: Notifies all waiting goroutines

### Key Implementation Pattern:

```go
// Creating a condition variable
cond := sync.NewCond(&sync.Mutex{})

// Waiting pattern
cond.L.Lock()
for !condition {
    cond.Wait()  // Automatically releases lock while waiting
}
// Do something with the condition met
cond.L.Unlock()

// Signaling pattern
cond.L.Lock()
// Change state
cond.Signal()  // or cond.Broadcast()
cond.L.Unlock()
```

### Example 1: Resource Management

```go
type Resource struct {
    n    int        // Number of available resources
    cond *sync.Cond // Condition variable for synchronization
}

func NewResource(n int) *Resource {
    r := new(Resource)
    r.n = n
    r.cond = sync.NewCond(&sync.Mutex{})
    return r
}

func (r *Resource) Acquire() {
    r.cond.L.Lock()
    for r.n == 0 {
        r.cond.Wait() // Wait until resources become available
    }
    r.n--
    r.cond.L.Unlock()
}

func (r *Resource) Release() {
    r.cond.L.Lock()
    r.n++
    r.cond.Signal() // Notify one waiting goroutine
    r.cond.L.Unlock()
}
```

### Example 2: Producer-Consumer Pattern

```go
type Buffer struct {
    capacity  int         // Maximum elements
    size      int         // Current elements
    mutex     sync.Mutex  // Protects the buffer
    notFull   *sync.Cond  // Signal when buffer is not full
    notEmpty  *sync.Cond  // Signal when buffer is not empty
}

func NewBuffer(capacity int) *Buffer {
    b := new(Buffer)
    b.capacity = capacity
    b.size = 0
    b.mutex = sync.Mutex{}
    b.notFull = sync.NewCond(&b.mutex)
    b.notEmpty = sync.NewCond(&b.mutex)
    return b
}

func (b *Buffer) Put() {
    b.mutex.Lock()
    for b.size == b.capacity {
        b.notFull.Wait() // Wait until buffer is not full
    }
    b.size++
    b.notEmpty.Signal() // Signal to waiting consumers
    b.mutex.Unlock()
}

func (b *Buffer) Get() {
    b.mutex.Lock()
    for b.size == 0 {
        b.notEmpty.Wait() // Wait until buffer is not empty
    }
    b.size--
    b.notFull.Signal() // Signal to waiting producers
    b.mutex.Unlock()
}
```

## 2. Channels

Channels provide a type-safe way for goroutines to communicate with each other.

Key properties:

- Send and receive data with the `<-` operator
- Direction of arrow indicates data flow
- Sends and receives block until other side is ready
- Nil channels block forever
- Closed channels can be read from but not written to

### Channel States:

- **Nil channel**: `var c chan int` — Blocks on send/receive
- **Open channel**: `c := make(chan int)` — Ready for data transfer
- **Closed channel**: `close(c)` — Can read remaining data, cannot write

### Example 1: Using Channels as Mutex

```go
func inc(c chan bool, wg *sync.WaitGroup) {
    for i := 0; i < 100_000; i++ {
        <-c        // Take token (lock)
        count++
        c <- true  // Return token (unlock)
    }
    wg.Done()
}

func main() {
    c := make(chan bool, 1) // Buffered channel simulates mutex
    var wg sync.WaitGroup

    wg.Add(2)
    go inc(c, &wg)
    go inc(c, &wg)

    c <- true  // Initialize with one token
    wg.Wait()
}
```

### Example 2: Select Statement with Multiple Channels

```go
func worker(c chan int) {
    // Generate random numbers periodically
    for {
        time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
        c <- rand.Intn(100)
    }
}

func main() {
    c1 := make(chan int)
    c2 := make(chan int)

    go worker(c1)
    go worker(c2)

    // Handle whichever channel has data ready
    for {
        select {
        case x1 := <-c1:
            fmt.Println("channel 1:", x1)
        case x2 := <-c2:
            fmt.Println("channel 2:", x2)
        }
    }
}
```

### Example 3: Handling Channel Closure

```go
func worker(c chan int) {
    for i := 0; i < 20; i++ {
        time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
        c <- rand.Intn(100)
    }
    close(c)  // Signal that no more values will be sent
}

func main() {
    c1 := make(chan int)
    c2 := make(chan int)

    go worker(c1)
    go worker(c2)

    for {
        select {
        case x1, ok := <-c1:
            if ok {
                fmt.Println("channel 1:", x1)
            } else {
                c1 = nil  // Set to nil to ignore this channel
            }
        case x2, ok := <-c2:
            if ok {
                fmt.Println("channel 2:", x2)
            } else {
                c2 = nil  // Set to nil to ignore this channel
            }
        }

        // Exit when both channels are closed
        if c1 == nil && c2 == nil {
            break
        }
    }
}
```

## 3. Stages

With channels, we can create a pipeline of stages, where each stage is a goroutine that processes data and sends it to the next stage.

```go

```

We can also make a structure like `Barrier` to synchronize multiple goroutines.

```go

```
