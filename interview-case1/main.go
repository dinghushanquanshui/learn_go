package main

import (
	"fmt"
)

func main() {
	// q1
	q1()

	// q2
	q2()

	// q3: <img src =xxx.png” alt =“免文件贷款”>
}

func q1() {
	tmp := []interface{}{1, 3, "a", "b", 4, 7}
	res := []int{}
	for _, v := range tmp {
		if num, ok := v.(int); ok {
			res = append(res, num)
		}
	}
	fmt.Println(res)
}

func q2() {
	tmp := []interface{}{5, "a", "b", "c", "d", "e", "f", "g"}
	firstVal := tmp[0].(int) + 1
	res := tmp[1:firstVal]
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		if i == 0 || i == len(res)-1 {
			continue
		}
		res[i], res[j] = res[j], res[i]
	}
	fmt.Println(res)
}
