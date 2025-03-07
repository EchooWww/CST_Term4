package main

import "fmt"
var a, b int

func main() {
	fmt.Println("a, b:", a, b)
	a, b := 10, 20
	fmt.Println("a, b:", a, b)
	arrays()
	slices()
}

func arrays() {
	fmt.Println("ARRAYS")
	var a [5]int
	fmt.Println("a:", a)
	a = [5]int{1, 2, 3, 4, 5}
	b := [5]int{5, 4, 3, 2, 1}
	c := [...]int{6, 7, 8, 9, 10}
	fmt.Println("a:", a)
	fmt.Println("b:", b)
	fmt.Println("c:", c)
	fmt.Println("sum5(a):", sum5(a))
	double(&a)
	fmt.Println("a:", a)
}

func sum_slices (s []int) int {
  sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

func slices() {
	var x[]int  // nil slice
	y := []int{} // empty slice
	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println(x == nil)
	fmt.Println(y == nil)
	a:= []int{1, 2, 3, 4, 5} 
	fmt.Println("a:", a, len(a), cap(a))
	b:= a[1:3]
	fmt.Println("b:", b, len(b), cap(b))
	x = append(x, 1) // append to a nil slice
	fmt.Println("x:", x)
}

func maps() {
	fmt.Println("MAPS")
	var x map[string]int // nil map
	y := map[string]int{} // empty map
	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println(x == nil)
	fmt.Println(y == nil)
	y["homer"] = 75000
	fmt.Println(y["homer"], y["marge"]) 
}

func sum5(a [5]int) int {
	sum := 0
	for i := 0; i < 5; i++ {
		sum += a[i]
	}
	return sum
}

func double(a *[5]int) {
	for i := 0; i < 5; i++ {
		(*a)[i] *= 2
	}
}