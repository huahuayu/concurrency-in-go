package main

import (
	"fmt"
	"sync"
)

func main(){
	var wg sync.WaitGroup
	for _, greeting := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(g string) {
			defer wg.Done()
			fmt.Println(g)
		}(greeting)
	}
	wg.Wait()
}


/*
output:
good day
greetings
hello

go func()匿名函数使用了传参的方式，这样三个goroutine各自有各自的
 */