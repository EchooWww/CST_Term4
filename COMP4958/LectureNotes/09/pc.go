// producer-consumer problem

package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)
type Buffer struct { // abstract to simulate just the capacity, size, free slots 
	capacity int // maximum number of elements the buffer can hold
	size int // number of elements currently in the buffer
	mutex sync.Mutex // to protect the buffer // we need the mutex explicitly as it is shared between the producer and consumer
	notFull, notEmpty *sync.Cond // condition variables to wait for the buffer to be not full or not empty. They should be the same mutex
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
	for b.size == b.capacity { // as long as buffer is full
		b.notFull.Wait() // wait for the buffer to be not full
	}
	fmt.Print("+")
	b.size++
	b.notEmpty.Signal()
}

func (b *Buffer) Get() {
	b.mutex.Lock()
	for b.size == 0 {
		b.notEmpty.Wait()
	}
	fmt.Print("-")
	b.size--
	b.notFull.Signal()
}

func produce(b *Buffer) {
	for {
		b.Put()
	}
}

func consume(b *Buffer) {
	for {
		b.Get()
	}
}

func main() {
	b := NewBuffer(100)
	go produce(b)
	go consume(b)
	go produce(b)
	go produce(b)
	go consume(b)

	_ = bufio.NewScanner(os.Stdin).Scan()
}