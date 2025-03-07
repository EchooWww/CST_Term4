package main

import "fmt"

type Circle struct {
	x, y, r int
}

func (c Circle) Draw() { // value receiver
	fmt.Printf("Draw a circle at (%d, %d) with radius %d\n", c.x, c.y, c.r)
}

func (c *Circle) Scale(factor int) { // pointer receiver
	c.r *= factor
}

func main() {
	c1 := Circle{1, 2, 3}
	c2 := Circle{
		x:1,
		y:2,
		r:3,
	}
	c1.Draw()
	c2.Draw()
}