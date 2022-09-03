package main


import "fmt"

func Increase() func() int {
    n := 0
    return func() int {
        n++
        return n
    }
}
func main() {
    i1 := Increase()
    i2 := Increase()
    fmt.Println(i1())
    fmt.Println(i2())
}