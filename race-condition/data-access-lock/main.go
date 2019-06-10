package main

import (
	"fmt"
	"sync"
)

func main(){
	var lock sync.Mutex

	var data int
	go func() {
		lock.Lock()
		data++
		lock.Unlock()
	}()

	lock.Lock()
	if data == 0 {
		fmt.Printf("the value is %v.\n", data)
	} else {
		fmt.Printf("the value is %v.\n", data)
	}
	lock.Unlock()
}

/*
1. 加锁解决了数据的竞争读写，go run --race main.go不会报错了
2. 但是没有解决竞争条件，程序执行路径仍是不确定的：
要么go routine先执行，要么if/else先执行，具体哪个先执行编程时并不知道。
*/
