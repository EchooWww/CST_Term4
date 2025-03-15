package main

type I interface {
	f() int
}

type S0 struct {}

func(s S0) f() int {
	return 0
} 

type S1 struct {}

func (s *S1) f() int {
	return 1
}

func main() {
	s0 := S0{}
	s1 := S1{}
	var _ I = s0
	var _ I = &s1
	// var _ I = s1 does not compile, because s1 is a pointer to S1
	// pointer receiver methods are more restrictive than value receiver methods
	var _ I = &s0
	// var _ *I = &s1 does not compile, because *I is not an interface, but I is an interface

}