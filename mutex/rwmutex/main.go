package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

func main(){
	//producer := func(wg *sync.WaitGroup, l sync.Locker) {
	//	defer wg.Done()
	//	for i := 5; i > 0; i-- {
	//		l.Lock()
	//		l.Unlock()
	//		time.Sleep(1)
	//	}
	//}
	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}
	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count)
		beginTestTime := time.Now()
		//go producer(&wg, mutex)             // what it is used for, sleep 5ns, if I comment it, the example still works
		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}
		wg.Wait()
		return time.Since(beginTestTime)
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()
	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\n")

	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &m, m.RLocker()),   // what's the difference between next line?
			test(count, &m, &m),
		)
	}
}
