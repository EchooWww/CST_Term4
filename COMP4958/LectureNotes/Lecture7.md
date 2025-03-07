# Lecture 7: Intro to GO

Everything in Go is a package. The main package
is the entry point of the program. The main function
is the entry point of the main package.

```go
package main
import "fmt" // format package

func main() {
    fmt.Println("Hello, World!")
}
```

If a member's name is capitalized, it is public.
If it is lowercase, it is private.

To run a Go program, use the `go run` command.

```bash
go run main.go
```

To build a Go program, use the `go build` command.

```bash
go build main.go
```

To import multiple

```go
import (
	"fmt"
	"os"
)
```

Anything uninitialized will contain the zero value

```go
var x int // 0
var c = 1 // type is inferred
d := "Hello" // type is inferred, short declaration
e, g := 1, 2 // multiple short declarations. But at least one must be new
fmt.Println(x, c, d) // can print multiple values

```

### Array

An array is a value instead of a reference. It is a fixed size.

```go
var a [5]int // array of 5 integers, containing 0s
a[2] = 7 // set the 3rd element to 7
b := [5]int{1, 2, 3, 4, 5} // array of 5 integers, initialized
c := [...]int{1, 2, 3, 4, 5} // array of 5 integers, initialized, size inferred
```

Go is a strict language: you cannot assign an array of size 5 to an array of size 6; you cannot declare a variable and not use it.

There's only one loop in Go: the `for` loop.

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

`i++` is a statement, not an expression. You cannot do `i = i++`.

### Pointers

The order of the type and the variable is reversed.

```go
func double (a *[5]int) {
    for i := range a {
        a[i] *= 2 // auto-dereference: a[i] is the same as (*a)[i]
    }
}
```

### Slices

Somewhat similar to rust slices: descriptor of contiguous segment of an array, just like a view.

A slice always has an underlying array.

But there's some complication in the Go version

Comparison between nil slice and empty slice: a nil slice doesn't have an underlying array, while an empty slice does, but it has a length of 0.

A slice has address + length + capacity. By default, the length and the capacity are the same, but capacity can be increased because of reallocation.

When we make a slice from an array, it is a view of the array. If we change the slice, we change the array.

```go
a := [5]int{1, 2, 3, 4, 5}
s := a[1:3] // slice of a, from 1 to 3
s[0] = 7 // a is now [1, 7, 3, 4, 5]
t := append(s, 8) // append 8 to s, t is now [7, 3, 8],
```

But when the slice exceeds the capacity, the slice is reallocated, but the value of the array is not changed.

We need slices when we want to make functions work for any size of array, but be careful when we append to or modify the slice.

Another way to use for loops

```go

func sum_slices (s []int) int {
  sum := 0
	for _, v := range s { // index + copy of value : if we modify v, the original value is not changed
		sum += v
	}
	return sum
}
b := [3]int{1, 2, 3}
fmt.Println(sum_slices(b[:])) // we need to first create a slice from the array
```

We can append to a nil slice

```go
	x = append(x, 1) // append to a nil slice
	fmt.Println("x:", x)
```

### Maps

We cannot use a nil map: go will panic.

If the key is not in the map, the value is the zero value of the type, just like in the case of an uninitialized variable.

```go
var m map[string]int // nil map
var n = map[string]int{} // empty map
n["one"] = 1
fmt.Println(n["one"], n["two"]) // 1, 0
```

Then how to test if a key is in the map?

```go
a, ok := n["two"] // indexing returns 2 values: value and a boolean
if ok {
    fmt.Println(a)
}
if b, ok := n["two"]; ok { // combine the two statements
    fmt.Println(b)
}
```

We can create a map using the `make` function

```go
o := make(map[string]int)
o["one"] = 1
z := map[string]int{"one": 1, "two": 2, "three": 3,} // the trailing comma is necessary
```

map is a reference type.

Loop through a map

```go
for k, v := range z {
    fmt.Println(k, ":", v)
}
```

### functions

Go is not a functional language: we can modify the arguments. But it still has some functional features: functions are first-class citizens.

We can make a function that returns a function

```go
func makeCounter(start int) func() int {
	return func() int {
		value := start
		start++
		return value
	}
}
```

We can also pass functions as arguments

```go

func find(s []int, p func(int) bool) (int, bool) {
	for _, x := range s {
		if p(x) {
			return x, true
		}
	}
	return 0, false
}
```

### Structs in Go
