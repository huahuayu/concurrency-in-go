package main

import (
	"fmt"
	"time"
)

func addEle(a []int, i int) {
	defer func() { //匿名函数捕获错误
		err := recover()
		if err != nil {
			fmt.Println("add ele fail", err)
		}
	}()
	a[i] = i
	fmt.Println(a)
}

func main() {
	is := make([]int, 4)
	for i := 0; i < 10; i++ {
		go addEle(is, i)
	}
	time.Sleep(time.Second * 2)
}
