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
