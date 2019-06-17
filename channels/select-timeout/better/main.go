package main

import (
	"fmt"
	"time"
)

func main() {
	var msg string
	ch := make(chan string, 1)
	defer close(ch)

	timer := time.NewTimer(1 * time.Microsecond)
	defer timer.Stop()

	go func() {
		//time.Sleep(1 * time.Microsecond) // uncomment to timeout
		ch <- "hi"
	}()

	select {
	case msg = <-ch:
		fmt.Println("Read from ch:", msg)
	case <-timer.C:
		fmt.Println("Timed out")
	}
}
