package main

import (
	"fmt"
	"time"
)

func send (c chan int) {
	time.Sleep(1 * time.Second)
	c <- 1
	fmt.Println("sent")
}

func main() {
	var c chan int // channel variable c of type int, nil channel
	// <-c // receive from nil channel, deadlock
	// c <- 1 // send to nil channel, deadlock
	c = make(chan int) // unbuffered channel: read and write must be synchronized. Read blocks until write, write blocks until read
	go send(c)
	// time.Sleep(2 * time.Second)
	x := <-c // wait for the value to be sent
	fmt.Println(x)

	_ := make(chan int, 10) // when you create a channel with a capacity, it is a buffered channel, we can send to it without waiting until the buffer is full
	close(c) // close the channel, no more values can be sent to it, but we can still receive values from it // it is always ready for reading 
}