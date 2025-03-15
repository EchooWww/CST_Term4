package main

import (
	"fmt"
	"sync"
	// "time" // mutex is in sync package
)

var mutex sync.Mutex // 0 value, unlocked
var count = 0
var wg sync.WaitGroup

func inc() { // not atomic: so it is the critical section needs to be protected
	defer wg.Done() // defer is executed when the function returns, no matter how it returns	
	for i := 0; i < 100_000; i++ {
		mutex.Lock()
		count++
		mutex.Unlock()
	}
	// wg.Done()
}

func main() {
	wg.Add(2)
	go inc()
	go inc()
 
	wg.Wait() // wait untio all goroutines are done
	// time.Sleep(2 * time.Second)
	fmt.Println(count) // 0. When main thread finishes, the program terminates, so the goroutines are not executed
}