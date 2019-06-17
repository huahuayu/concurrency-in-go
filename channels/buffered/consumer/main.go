package main

import (
	"fmt"
	"sync"
)

func Producer(mychan chan int, data int, wait *sync.WaitGroup) {
	mychan <- data
	fmt.Println("produce data：", data)
	wait.Done()
}
func Consumer(mychan chan int, wait *sync.WaitGroup) {
	a := <-mychan
	fmt.Println("consume data：", a)
	wait.Done()
}
func main() {

	datachan := make(chan int, 5) //通讯数据管道
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		go Producer(datachan, i, &wg) //生产数据
		wg.Add(1)
	}
	for j := 0; j < 10; j++ {
		go Consumer(datachan, &wg) //消费数据
		wg.Add(1)
	}
	wg.Wait()
}
