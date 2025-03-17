package main

import (
	"bug-free/lec9/util"
	"fmt"
	"math/rand"
)

func worker(c chan int) {
	for i:=0; i < 20; i++{
		util.Work(500, 1000)
		c <- rand.Intn(100)
	}
	close(c) // close is not necessary unless we want to signal that no more values will be sent on the channel
}

func main() {
	c1 := make(chan int)
	c2 := make(chan int)

	go worker(c1)
	go worker(c2)

	for {
		select { // select statement allows us to wait on multiple channels
		case x1, ok := <- c1:
			if ok {
				fmt.Println("channel 1:", x1)
			} else {
				c1 = nil
			}
		case x2, ok := <- c2:
			if ok {
				fmt.Println("channel 2:", x2)
			} else {
				c2 = nil
			}
		}
		if c1 == nil && c2 == nil {
			break
		}
	}
}