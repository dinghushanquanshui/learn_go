package main

import "fmt"

type MyInterface interface {
	M1()
}

type T int

func (T) M1() {
	println("T's M1")
}

// 如果 v,ok := i.(T)，T 是一个接口类型，意味着断言 i 的值实现了接口类型 T；
// 如果断言成功，变量 v 的类型为 i 的值的类型；如果断言失败，v 的类型信息为接口类型 T，它的值为 nil
func main() {
	var t T
	var i interface{} = t
	v1, ok := i.(MyInterface)
	if !ok {
		panic("the value of i is not MyInterface")
	}

	v1.M1()
	fmt.Printf("the type of  v1 is %T\n", v1)

	i = int64(13)
	v2, ok := i.(MyInterface)
	fmt.Printf("the type of v2 is %T\n", v2)
	// v2 = 13
}
