package main

import (
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
