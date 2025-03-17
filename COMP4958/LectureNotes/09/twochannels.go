package main

import (
	"bug-free/lec9/util"
	"fmt"
	"math/rand"
)

func worker(c chan int) {
	for {
		util.Work(500, 1000)
		c <- rand.Intn(100)
	}
}

func main() {
	c1 := make(chan int)
	c2 := make(chan int)

	go worker(c1)
	go worker(c2)

	for {
		select { // select statement allows us to wait on multiple channels
		case x1 := <- c1:
			fmt.Println("channel 1:", x1)
		case x2 := <- c2:
			fmt.Println("channel 2:", x2)
		}
	}
}