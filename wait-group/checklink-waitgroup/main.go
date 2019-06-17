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

func checkLink(wg *sync.WaitGroup, link string) {
	defer wg.Done()
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "might be down! ", err)
		return
	}
	fmt.Println(link, "is up!")
}
