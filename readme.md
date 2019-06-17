# concurrency in go
go关键字 goroutines/example1  
goroutine与匿名函数 goroutines/example2    
函数一等公民  goroutines/example3  
竞态条件 race-condition/simplest-race-condition  
数据锁 race-condition/data-access-lock  
goroutine内存消耗测试  goroutine-benchmark  
闭包 goroutine-and-closure/simple-closure  
闭包 goroutine-and-closure/fibonacci  
闭包和goroutine常见错误 goroutine-and-closure/common-mistake  
闭包和goroutine常见错误避免 goroutine-and-closure/common-mistake-avoid  
错误恢复 recover   
waitgroup helloworld wait-group/helloworld-waitgroup  
waitgroup一般用法 wait-group/waitgroup-common-pattern  
互斥锁 mutex/mutex  
死锁 dead-lock  
活锁 live-lock  
饥饿 live-lock-starvation 
通道 channels/unbuffered  
带缓冲的通道 channels/buffered  
不安全的工作池 channels/simpleworkpoll/unsafe  
安全的工作池 channels/simpleworkpoll/safe  
select channels/select  
select + 工作池 channels/selectjobqueue  
通道的方向 channels/channeldirection  
上下文 channels/context  
带value的上下文 channels/contextwithvalue  
多协程的上下文 channels/mutigoroutincontext    
综合示例：  
模拟google搜索 googlesearch  





## 并发和并行
并发（Concurrency）vs 并行（Parallelism）的区别  
你吃饭吃到一半，电话来了，你一直到吃完了以后才去接，这就说明你不支持并发也不支持并行。  
你吃饭吃到一半，电话来了，你停了下来接了电话，接完后继续吃饭，这说明你支持并发。  
你吃饭吃到一半，电话来了，你一边打电话一边吃饭，这说明你支持并行。  
![](https://i.imgur.com/us17QJ2.jpg)

golang同时支持并行和并发，并行即使用多个cpu同时运算（默认使用所有cpu，可由runtime.GOMAXPROCS(n int)指定)  
``` go
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("i: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
```

## race condition
开发并发程序很容易出错，因为开发者容易陷入一种顺序的思考，错误的认为先写的代码就会先执行，这在写并发程序的时候是不成立的。  

以下程序运行可能有三种结果：  
1. 什么也不打印，第11行先于14行执行，程序运行结束  
2. "the value is 0"被打印，第14和15行先于第11行执行  
3. "the value is 1"被打印，第14行先与11行执行，但是15行在11行后执行  

``` go
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

```

加sleep的方式不可取，这样做没有根本性解决问题，race condition依然存在，而且降低系统效率  
``` go
package main

import (
	"fmt"
)

func main(){
	var data int

	go func(){
		data++
	}()

    time.Sleep(1 * time.Second) // 错误示范
    
	if data == 0{
		fmt.Printf("the value is %v.\n",data)
	}
}
```

## 检查race condition
使用`--race`参数可以检查程序是否有race condition  
``` go
go run --race main.go 
```

## channel
> Do not communicate by sharing memory; instead, share memory by communicating.

