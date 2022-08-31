package main

import "fmt"

type T int

func (t T) Error() string {
	return "bad error"
}

func printNonEmptyInterface() {
	var err1 error
	var err2 error
	err1 = (*T)(nil)
	println("err1:", err1)
	println("err1 == nil", err1 == nil)

	err1 = T(5)
	err2 = T(6)
	println("err1:", err1)
	println("err2:", err2)
	println("err1 == err2:", err1 == err2)

	err2 = fmt.Errorf("%d\n", 5)
	println("err1:", err1)
	println("err2:", err2)
	println("err1 == err2:", err1 == err2)
}

func main() {
	//
	printNonEmptyInterface()
}
