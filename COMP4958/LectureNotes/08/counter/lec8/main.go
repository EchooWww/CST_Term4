package main

import (
	"fmt"
	"lec8/util"
	"sync"
)

func inc(c *util.Counter, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100_000; i++ {
		c.Inc()
	}
}


func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	c := util.NewCounter(0)
	go inc(c, &wg)
	go inc(c, &wg)
	wg.Wait()
	fmt.Println(c.Value)
}