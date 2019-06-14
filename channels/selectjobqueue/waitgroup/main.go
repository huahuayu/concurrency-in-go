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

func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())

	// make a channel with a capacity of 10.
	jobChan := make(chan Job, 1000)

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
	WaitTimeout(&wg, 2*time.Second)
	elapsed := time.Since(start)
	fmt.Printf("%d work assigned, %d been done, in %s", workAssign, workDone, elapsed)
}
