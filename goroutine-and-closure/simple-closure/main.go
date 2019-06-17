package main

import (
	"fmt"
)

func main() {
	greeting := "hello"

	go func() {
		greeting = "welcome"
	}()

	fmt.Println(greeting)
}
