package main

import (
	"fmt"
	"time"
)

func main(){
	go SayHello()

	fmt.Println("hello from main")
	time.Sleep(time.Second)
}

func SayHello(){
	fmt.Println("hello from goroutine")
}
