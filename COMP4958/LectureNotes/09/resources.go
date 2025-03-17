package main

import (
	"bug-free/lec9/util"
	"fmt"
	"sync"
)

type Resource struct {
	// mutex sync.Mutex // we don't need the mutex when we have a condition variable
	n int // number of resources maximum available
	cond *sync.Cond // we need a condition variable to wait for resources
}

func NewResource(n int) *Resource {
	r := new(Resource)
	r.n = n
	r.cond = sync.NewCond(&sync.Mutex{})
	return r // the mutex is initialized to the 0 value
}

// 2 methods to acquire and release resources
func (r *Resource) Acquire() {
	r.cond.L.Lock()
	for r.n == 0 {
		r.cond.Wait()
	}
	fmt.Printf("- %d -> %d\n", r.n, r.n-1)
	r.n--
	r.cond.L.Unlock()
}

func (r *Resource) Release() {
	r.cond.L.Lock()
	fmt.Printf("+ %d -> %d\n", r.n, r.n+1)
	r.n++
	r.cond.Signal() // let other goroutines know that a resource is available
	r.cond.L.Unlock()
}

func user(r *Resource) {
	for i := 0; i < 10; i++ {
		util.Work(500, 1000)
		r.Acquire()
		util.Work(500, 1000)
		r.Release()
	} 
}

var wg sync.WaitGroup

func main() {
	r := NewResource(3)
	wg.Add(5)
	go user(r)
	go user(r)
	go user(r)
	go user(r)
	go user(r)
	wg.Wait()
}