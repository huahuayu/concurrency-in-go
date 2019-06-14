package gpool

import (
	"sync"
)

type pool struct {
	queue chan int
	wg    *sync.WaitGroup
}

func New(size int) *pool {
	if size <= 0 {
		size = 1
	}
	return &pool{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}

func (p *pool) Add(delta int) {
	for i := 0; i < delta; i++ {
		p.queue <- 1
	}
	for i := 0; i > delta; i-- {
		<-p.queue
	}
	p.wg.Add(delta)
}

func (p *pool) Done() {
	<-p.queue
	p.wg.Done()
}

func (p *pool) Wait() {
	p.wg.Wait()
}

/*
usage:
	pool := gpool.New(5)
	for _,l := range links {
		pool.Add(1)
		go func(link string) {
			checkLink(link)
			time.Sleep(time.Second)
			pool.Done()
		}(l)
	}

	pool.Wait()
*/
