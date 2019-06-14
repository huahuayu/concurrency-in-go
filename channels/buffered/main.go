package main

import (
	"GoCasts/code/channels/buffered/gpool"
	"fmt"
	"net/http"
	"time"
)

func main() {
	links := []string{
		"http://baidu.com",
		"http://qq.com",
		"http://taobao.com",
		"http://jd.com",
		"http://z.cn",
	}

	pool := gpool.New(5)
	for {
		time.Sleep(2 * time.Second)
		for _, l := range links {
			pool.Add(1)
			go func(link string) {
				checkLink(link)
				pool.Done()
			}(l)
		}
	}

	pool.Wait()
}

func checkLink(link string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "might be down!")
		return
	}

	fmt.Println(link, "is up!")
}
