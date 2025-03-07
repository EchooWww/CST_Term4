package main

import "fmt"

func main() {
	c1 := makeCounter(5)
	c2 := makeCounter(10)
	fmt.Println(c1())
	fmt.Println(c1())
	fmt.Println(c2())
	fmt.Println(c2())

	a := [...]int{3,2,7,6,8}
	if x, ok := find(a[:], func(n int) bool {return n % 6 == 0}); ok {
		fmt.Println(x)
	} else {
		fmt.Println("Not found")
	}
}

func makeCounter(start int) func() int {
	return func() int {
		value := start
		start++
		return value
	}
}

func find(s []int, p func(int) bool) (int, bool) {
	for _, x := range s {
		if p(x) {
			return x, true
		}
	}
	return 0, false
}