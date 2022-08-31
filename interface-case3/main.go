package main

func main() {
	// nil interface
	var i interface{}
	var err error

	// 无论是空接口类型还是非空接口类型，一旦变量值为 nil, 那么它们内部表示均为 (0x0,0x0)

	println(i)
	println(err)
	println("i = nil:", i == nil)
	println("err = nil:", err == nil)
	println("i = err:", i == err)
}
