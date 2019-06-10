package main

import (
	"fmt"
)

func main(){
	var data int

	go func(){
		data++
	}()

	if data == 0{
		fmt.Printf("the value is %v.\n",data)
	}
}
/*
程序运行可能有三种结果：
1. 什么也不打印，第11行先于14行执行，程序运行结束
2. "the value is 0"被打印，第14和15行先于第11行执行
3. "the value is 1"被打印，第14行先与11行执行，但是15行在11行后执行
*/