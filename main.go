package main

import "fmt"

func test(a int) int {
	for true {

	}
}

const (
	a = 1 << iota
	b
	c
)

type S[T any] struct {
	t T
}

func main() {

	x := S[]{}

	var f int
	f |= a
	fmt.Println(f)

	f |= b
	fmt.Println(f)

	f |= c
	fmt.Println(f)

	f = f & ^c
	fmt.Println(f)

	f = f & ^b
	fmt.Println(f)

	f = f & ^a
	fmt.Println(f)
}
