package main

import (
	"fmt"
	"time"
)

func main() {
	for _, greeting := range []string{"hello", "greetings", "good day"} {
		go func(words string) {
			fmt.Println(words)
		}(greeting)
	}

	time.Sleep(1 * time.Second)
}

/*
output:
good day
greetings
hello

go func()匿名函数使用了传参的方式，这样三个goroutine各自有各自的本地变量
*/
