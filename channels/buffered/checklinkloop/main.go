package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	links := []string{
		"http://qq.com",
		"http://taobao.com",
		"http://baidu.com",
		"http://z.cn",
		"http://great.website",
	}

	c := make(chan string)

	for _, link := range links {
		go checkLink(link, c)
	}

	for l := range c {
		go func(link string) {
			time.Sleep(5 * time.Second)
			checkLink(link, c)
		}(l)
	}
}

func checkLink(link string, c chan string) { // could be any routine
	defer func() {c <- link}()

	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "might be down!" )
		return
	}

	fmt.Println(link, "is up!")
}
