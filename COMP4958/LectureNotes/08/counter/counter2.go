package main

import (
	"fmt"
	"sync"
	// "time" // mutex is in sync package
)

type Counter struct {
	mutex sync.Mutex
	Value int
}

func NewCounter(value int) *Counter {
	c := new(Counter)
	c.Value = value
	c.mutex = sync.Mutex{}
	return c
}

func (c *Counter) Inc() {
	c.mutex.Lock()
	c.Value++
	c.mutex.Unlock()
}

func inc(c *Counter, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100_000; i++ {
		c.Inc()
	}
}


func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	c := NewCounter(0)
	go inc(c, &wg)
	go inc(c, &wg)
	fmt.Println(c.Value)
}