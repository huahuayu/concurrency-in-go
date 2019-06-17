package main

import (
	"fmt"
	"time"
)

func main() {
	for _, greeting := range []string{"hello", "greetings", "good day"} {
		go func() {
			fmt.Println(greeting)
		}()
	}

	time.Sleep(1 * time.Second)
}

/*
output:
good day
good day
good day

因为使用了go关键字，for循环里面的逻辑执行是非阻碍的，for循环瞬间执行完了，在三个go routine执行fmt.Println时，
三个goroutine使用的都是同一个变量地址空间，取到的都是同一个值，"good day"
*/
