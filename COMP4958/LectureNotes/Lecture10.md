# Lecture 10

Doing seive of eratosthenes with gorooutines: we can have a goroutine for each of the filters. We are also using channels to simulate a stream of data.

```go
package main

import (
  "fmt"
)

func generate(c chan int) {
  for i := 2; ; i++ {
    c <- i
  }
}

func filter(p int, in, out chan int) {
  for {
    i := <-in
    if i % p != 0 {
      out <- i
    }
  }
}

func main() {
  in := make(chan int)
  go generate(in) // create a goroutine to generate numbers
  for i := 0; i < 100; i++> { // first 100 prime numbers
    p := <-in
    fmt.Println(p)
    out := make(chan int)
    go filter(p, in, out) // create a goroutine to filter out multiples of p
    in = out // this is the new input channel
  }
}
```
