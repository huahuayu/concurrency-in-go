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
		wg.Wait()
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
