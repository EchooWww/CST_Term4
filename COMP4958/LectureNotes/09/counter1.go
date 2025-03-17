package main

import (
	"fmt"
	"sync"
)

var count = 0

func inc(c chan bool) { // we can use a channel to signal the completion of the goroutine
	for i := 0; i < 100_000; i++ {
		<-c
		count++
		c <- true
	}
	wg.Done()
}

var wg sync.WaitGroup
func main() {
	c := make(chan bool, 1) // buffered channel with capacity 1 can simulate a mutex: we need the size 1 buffer because there sould be an extra write to the channel
	wg.Add(2)
	go inc(c)
	go inc(c)
	c <- true
	wg.Wait()
	// time.Sleep(1 * time.Second)
	fmt.Println(count) // 200_000
}