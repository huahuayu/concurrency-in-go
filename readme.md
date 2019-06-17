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
带方向的channel channels/channeldirection  
带缓冲的通道 channels/buffered  
简单工作队列 channels/simplejobqueue  
关闭通道（错误示范） channels/simpleworkpoll/unsafe  
关闭通道（正确示范） channels/simpleworkpoll/safe   
非阻塞channel(select) channels/select  
通道timeout channels/select-timeout   
select + 工作池 channels/selectjobqueue  
通道的方向 channels/channeldirection  
上下文 channels/context  
带value的上下文 channels/contextwithvalue  
多协程的上下文 channels/mutigoroutincontext    
综合示例：  
模拟google搜索 googlesearch  





## 并发 vs 并行
并发（Concurrency）vs 并行（Parallelism）的区别  
你吃饭吃到一半，电话来了，你一直到吃完了以后才去接，说明你不支持并发也不支持并行。  
你吃饭吃到一半，电话来了，你停了下来接了电话，接完后继续吃饭，这说明你支持并发。  
你吃饭吃到一半，电话来了，你一边打电话一边吃饭，这说明你支持并行。  
![](https://i.imgur.com/us17QJ2.jpg)

golang同时支持并行和并发，并行即使用多个cpu同时运算（默认使用所有cpu，可由runtime.GOMAXPROCS(n int)指定)  
``` go
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("i: ", i)
		}()
	}
    time.Sleep(1 * time.Second)
}
```

## goroutine
只需使用go关键字即可启动一个goroutine，它使用异步非阻塞的方式运行。
``` go
// goroutines/example1 
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
```

## goroutine的消耗
``` go
// goroutine-benchmark/goroutine-memory
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c } // a goroutine which will never exit, because it's waiting for channel c all the time
	const numGoroutines = 1e5
	wg.Add(numGoroutines)
	before := memConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1024)
}
```

## 匿名函数
``` go
// goroutines/example2 
package main

import (
	"fmt"
)

func main(){
	go func(){
		fmt.Println("hello from goroutine")
	}()

	fmt.Println("hello from main")
}
```

## 函数一等公民
``` go
// goroutines/example3  
package main

import (
	"fmt"
)

func main() {
	sayHello := func() {
		fmt.Println("hello from goroutine")
	}

	go sayHello()

	fmt.Println("hello from main")
}
```

## 闭包
``` go
// goroutine-and-closure/simple-closure  
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
```

## 闭包 - 斐波那契数列
``` go
package main

import (
	"fmt"
)

func main() {
	fibonacci := func() func() int {
		back1, back2 := -1, 1
		return func() int {
			back1, back2 = back2, (back1 + back2)
			return back2
		}
	}

	f := fibonacci()

	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
```

## 闭包常见错误
``` go
package main

import (
	"fmt"
	"time"
)

func main() {
	for _, greeting := range []string{"hello", "greetings", "good day"} {
		go func() {
			fmt.Println(greeting)
		}()
	}

	time.Sleep(1 * time.Second)
}
```

## 闭包常见错误避免
``` go
package main

import (
	"fmt"
	"time"
)

func main() {
	for _, greeting := range []string{"hello", "greetings", "good day"} {
		go func(words string) {
			fmt.Println(words)
		}(greeting)
	}

	time.Sleep(1 * time.Second)
}
```

## 竞态条件
以下程序运行有几种结果？
``` go
// race-condition/simplest-race-condition 
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
程序运行可能有三种结果：
1. 什么也不打印，第11行（data++）先于14行执行，程序运行结束
2. "the value is 0"被打印，第14和15行先于第11行执行
3. "the value is 1"被打印，第14行先与11行执行，但是15行在11行后执行

## 互斥锁
``` go
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
```


## 示例1:网站健康检查
功能：检查一批网站是否能正常访问  
``` go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	links := []string{
		"http://qq.com",
		//"http://google.com",
		"http://taobao.com",
		"http://baidu.com",
		"http://z.cn",
		"http://great.website",
	}

	for _, link := range links {
		checkLink(link)
	}
}

func checkLink(link string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "might be down!", err)
		return
	}
	fmt.Println(link, "is up!")
}
```

问题：如果有部分网站访问超时（把google.com的注释打开），程序会怎么样？怎么避免这样的情况？  



## 示例2:网站健康检查（goroutine）
功能：同时检查一批网站是否能正常访问  
``` go
package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	links := []string{
		"http://qq.com",
		"http://google.com",
		"http://taobao.com",
		"http://baidu.com",
		"http://z.cn",
		"http://great.website",
	}

	for _, link := range links {
		go checkLink(link)
	}

	time.Sleep(5 * time.Second)
}

func checkLink(link string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "might be down! ", err)
		return
	}
	fmt.Println(link, "is up!")
}
```
问题1：输出的结果是什么？
问题2：如果把sleep拿走，结果会怎么样？为什么？

## waitgroup
``` go
// wait-group/helloworld-waitgroup 
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}
	wg.Wait()
	
	fmt.Println("All goroutines complete.")
}
```

## 示例3：网站健康检查（waitgroup）
``` go
package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	links := []string{
		"http://qq.com",
		//"http://google.com",
		"http://taobao.com",
		"http://baidu.com",
		"http://z.cn",
		"http://great.website",
	}

	for _, link := range links {
		wg.Add(1)
		go checkLink(&wg, link)
	}

	wg.Wait()

}
```

## channel
``` go
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
```

## 示例4：网站健康检查(channel)

``` go
// channels/unbuffered
package main

import (
	"fmt"
	"net/http"
	"time"
)

// 检查网站是否正常
func main() {
	links := []string{
		"http://baidu.com",
		"http://qq.com",
		"http://taobao.com",
		"http://jd.com",
		"http://z.cn",
	}

	c := make(chan string)

	for _, link := range links {
		go func(link string) {
			checkLink(link, c)
		}(link)
	}

	for l := range c {
		go func(link string) {
			time.Sleep(5 * time.Second)
			checkLink(link, c)
		}(l)
	}
}

func checkLink(link string, c chan string) {
	defer func() { c <- link }()

	_, err := http.Get(link)

	if err != nil {
		fmt.Println(link, "might be down!")
	} else {
		fmt.Println(link, "is up!")
	}

}
```

## 带方向的channel  
``` go
// channels/channeldirection
package main

import "fmt"

// This `ping` function only accepts a channel for sending
// values. It would be a compile-time error to try to
// receive on this channel.
func ping(pings chan<- string, msg string) {
	pings <- msg
}

// The `pong` function accepts one channel for receives
// (`pings`) and a second for sends (`pongs`).
func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
	// pongs <- <-pings
}

func main() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}
```

## 带缓冲区的channel  
``` go
// channels/simplejobqueue
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Job struct {
	JobName string
}

func worker(jobChan <-chan Job) {
	for job := range jobChan {
		process(job)
	}
}

func process(j Job) {
	fmt.Println(j.JobName + " processed")
}

func main() {
	// make a channel with a capacity of 100.
	jobChan := make(chan Job, 100)

	// start the worker
	go worker(jobChan)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	// enqueue a job
	for {
		i := r.Intn(10000)
		time.Sleep(time.Second)
		fmt.Printf("job %d assigned\n", i)
		job := Job{JobName: fmt.Sprintf("job " + strconv.Itoa(i))}
		jobChan <- job
	}

}
```

## 关闭通道
``` go
// channels/simpleworkpoll/unsafe 
package main

import (
	"fmt"
)

// Here's the worker, of which we'll run several
// concurrent instances. These workers will receive
// work on the `jobs` channel and send the corresponding
// results on `results`. We'll sleep a second per job to
// simulate an expensive task.
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func main() {

	// In order to use our pool of workers we need to send
	// them work and collect their results. We make 2
	// channels for this.
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for i := 0; i < 3; i++ {
		go worker(i, jobs, results)
	}

	// Here we send 5 `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for j := 0; j < 5; j++ {
		jobs <- j
	}
	close(jobs)

	// warning: it's not a good way to get the result
	for r := 0; r < 5; r++ {
		fmt.Println(<-results)
	}

}
```

## 关闭/range通道（正确示范）
``` go
// channels/simpleworkpoll/safe
package main

import (
	"fmt"
	"sync"
)

// Here's the worker, of which we'll run several
// concurrent instances. These workers will receive
// work on the `jobs` channel and send the corresponding
// results on `results`. We'll sleep a second per job to
// simulate an expensive task.
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func main() {

	// In order to use our pool of workers we need to send
	// them work and collect their results. We make 2
	// channels for this.
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	wg := new(sync.WaitGroup)

	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(i, jobs, results, wg)
	}

	// Here we send 5 `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for j := 0; j < 5; j++ {
		jobs <- j
	}
	close(jobs)

	go func(wg *sync.WaitGroup, results chan int) {
		fmt.Println("waiting")
		wg.Wait() // GOOD
		fmt.Println("done waiting")
		close(results)
	}(wg, results)

	//for r := 0; r < 5; r++ {
	//	fmt.Println(<-results)
	//}
	// Finally we collect all the results of the work.
	for r := range results {
		fmt.Println(r)
	}
}
```

## 非阻塞channel - select
``` go
// channels/select 
package main

import (
	"fmt"
	"time"
)

func main() {
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("监控退出，停止了...")
				return
			default:
				fmt.Println("goroutine监控中...")
				time.Sleep(2 * time.Second)
			}
		}
	}()

	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	stop <- true
	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)

}
```

## timeout
``` go
// channels/select-timeout
package main

import (
	"fmt"
	"time"
)

func main() {
	var msg string
	ch := make(chan string, 1)
	defer close(ch)

	go func() {
		//time.Sleep(1 * time.Microsecond)   // uncomment to timeout
		ch <- "hi"
	}()

	select {
	case msg = <-ch:
		fmt.Println("Read from ch:", msg)
	case <-time.After(1 * time.Microsecond):
		fmt.Println("Timed out")
	}
}
```

## 过载保护
``` go
// channels/selectjobqueue/waitgroup
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup
var workDone int
var workAssign int

type Job struct {
	JobName string
}

func worker(i int, jobChan <-chan Job) {
	defer wg.Done()
	for job := range jobChan {
		process(job)
		fmt.Println("worker: " + strconv.Itoa(i) + " " + job.JobName + " processed")
		workDone++
	}
}

func process(j Job) {
	// work to process
}

// TryEnqueue tries to enqueue a job to the given job channel. Returns true if
// the operation was successful, and false if enqueuing would not have been
// possible without blocking. Job is not enqueued in the latter case.
func TryEnqueue(job Job, jobChan chan<- Job) bool {
	select {
	case jobChan <- job:
		return true
	default:
		return false
	}
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())

	// make a channel with a capacity of 10.
	jobChan := make(chan Job, 10)

	// start the worker
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			worker(i, jobChan)
		}(i)

	}

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	// enqueue a job
	start := time.Now()
	for {
		i := r.Intn(10000)
		job := Job{JobName: fmt.Sprintf("job " + strconv.Itoa(i))}
		if !TryEnqueue(job, jobChan) {
			fmt.Println("max capacity reached, try later")
			close(jobChan)
			break
		} else {
			fmt.Printf("job %d assigned\n", i)
			workAssign++
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("%d work assigned, %d been done, in %s", workAssign, workDone, elapsed)
}
```

## 示例5：google search 1.0
![](https://raw.githubusercontent.com/huahuayu/img/master/20190616174410.png)
``` go
// googlesearch/googlesearch1.0
/*
Example: Google Search
Given a query, return a page of search results (and some ads).
Send the query to web search, image search, YouTube, Maps, News, etc. then mix the results.
Google function takes a query and returns a slice of Results (which are just strings)
Google invokes Web, Image and Video searches serially, appending them to the results slice.
Run each search in series
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	web   = fakeSearch("web")
	image = fakeSearch("image")
	video = fakeSearch("video")
)

type (
	result string
	search func(query string) result
)

func main() {
	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	results := google("golang")
	elapsed := time.Since(start)

	fmt.Println(results)
	fmt.Println(elapsed)
}

func fakeSearch(kind string) search {
	return func(query string) result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

func google(query string) (results []result) {
	results = append(results, web(query))
	results = append(results, image(query))
	results = append(results, video(query))

	return results
}
```
## 示例6：google search 2.0 
![](https://raw.githubusercontent.com/huahuayu/img/master/20190616174502.png)

``` go
// googlesearch/googlesearch2.0
/*
Example: Google Search
Given a query, return a page of search results (and some ads).
Send the query to web search, image search, YouTube, Maps, News, etc. then mix the results.
Run the Web, Image and Video searches concurrently, and wait for all results.
No locks. No condition variables. No callbacks
Run each search in their own Goroutine and wait for all searches to complete before display results
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	web   = fakeSearch("web")
	image = fakeSearch("image")
	video = fakeSearch("video")
)

type (
	result string
	search func(query string) result
)

func main() {
	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	results := google("golang")
	elapsed := time.Since(start)

	fmt.Println(results)
	fmt.Println(elapsed)
}

func fakeSearch(kind string) search {
	return func(query string) result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}
func google(query string) (results []result) {
	c := make(chan result)

	go func() {
		c <- web(query)
	}()

	go func() {
		c <- image(query)
	}()

	go func() {
		c <- video(query)
	}()

	for i := 0; i < 3; i++ {
		r := <-c
		results = append(results, r)
	}

	return results
}
```

## 示例7：google search 2.1
![](https://raw.githubusercontent.com/huahuayu/img/master/20190616174530.png)
``` go
// googlesearch/googlesearch2.1
/*
Example: Google Search 2.1
Given a query, return a page of search results (and some ads).
Send the query to web search, image search, YouTube, Maps, News, etc. then mix the results.
Don't wait for slow servers. No locks. No condition variables. No callbacks
Run each search in their own Goroutine but only return any searches that complete in
80 Milliseconds or less
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	web   = fakeSearch("web")
	image = fakeSearch("image")
	video = fakeSearch("video")
)

type (
	result string
	search func(query string) result
)

func main() {
	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	results := google("golang")
	elapsed := time.Since(start)

	fmt.Println(results)
	fmt.Println(elapsed)
}

func fakeSearch(kind string) search {
	return func(query string) result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}
func google(query string) (results []result) {
	c := make(chan result)

	go func() {
		c <- web(query)
	}()

	go func() {
		c <- image(query)
	}()

	go func() {
		c <- video(query)
	}()
	timeout := time.After(80 * time.Millisecond)

	for i := 0; i < 3; i++ {
		select {
		case r := <-c:
			results = append(results, r)
		case <-timeout:
			fmt.Println("timed out")
			return results
		}
	}

	return results
}

```


## 示例8：google search 3.0
![](https://raw.githubusercontent.com/huahuayu/img/master/20190616174602.png)
``` go
// googlesearch/googlesearch3.0
/*
Example: Google Search 3.0
Given a query, return a page of search results (and some ads).
Send the query to web search, image search, YouTube, Maps, News, etc. then mix the results.
No locks. No condition variables. No callbacks
Reduce tail latency using replicated search servers
Run the same search against multiple servers in their own Goroutine but only return searches
that complete in 80 Milliseconds or less
All three searches SHOULD always come back in under 80 milliseconds
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	web1   = fakeSearch("web")
	web2   = fakeSearch("web")
	image1 = fakeSearch("image")
	image2 = fakeSearch("image")
	video1 = fakeSearch("video")
	video2 = fakeSearch("video")
)

type (
	result string
	search func(query string) result
)

func main() {
	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	results := google("golang")
	elapsed := time.Since(start)

	fmt.Println(results)
	fmt.Println(elapsed)
}

func fakeSearch(kind string) search {
	return func(query string) result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}
func google(query string) (results []result) {
	c := make(chan result)

	go func() {
		c <- first(query, web1, web2)
	}()

	go func() {
		c <- first(query, image1, image2)
	}()

	go func() {
		c <- first(query, video1, video2)
	}()
	timeout := time.After(80 * time.Millisecond)

	for i := 0; i < 3; i++ {
		select {
		case r := <-c:
			results = append(results, r)
		case <-timeout:
			fmt.Println("timed out")
			return results
		}
	}

	return results
}

func first(query string, replicas ...search) result {
	c := make(chan result)

	// Define a function that takes the index to the replica function to use.
	// Then it executes that function writing the results to the channel.
	searchReplica := func(i int) {
		c <- replicas[i](query)
	}

	// Run each replica function in its own Goroutine.
	for i := range replicas {
		go searchReplica(i)
	}

	// As soon as one of the replica functions write a result, return.
	return <-c
}
```
