package main

import (
	"time"
)

func main() {
	var ch chan string     // 声明
	ch = make(chan string) // 初始化

	go ask(ch)
	go answer(ch)

	time.Sleep(1 * time.Second)
}

func ask(ch chan string) {
	ch <- "what's your name?"
}

func answer(ch chan string) {
	println("he asked: ", <-ch)
	println("My name is Shiming")
}
