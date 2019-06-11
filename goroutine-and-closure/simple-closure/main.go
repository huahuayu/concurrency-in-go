package main

import (
	"fmt"
	"sync"
)

func main(){
	var wg sync.WaitGroup
	greeting := "hello"

	wg.Add(1)
	go func() {    // 1. closure in a goroutine
		defer wg.Done()
		greeting = "welcome"   // 2. closure change outside variable or his own copy?
	}()
	wg.Wait()

	fmt.Println(greeting)   // 3. output: welcome, so the closure changed the outside variable
}
