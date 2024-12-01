package main

import "fmt"

//go:generate gogen stack.gogen int

func main() {
	is := intStack{}
	is.Push(1)
	is.Push(2)
	is.Push(3)
	fmt.Println(is.Pop())
	fmt.Println(is.Pop())
	fmt.Println(is.Pop())
}
