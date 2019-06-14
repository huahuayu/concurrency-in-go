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
	// make a channel with a capacity of 10.
	jobChan := make(chan Job, 10)

	// start the worker
	go worker(jobChan)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	// enqueue a job
	for {
		i := r.Intn(10000)
		job := Job{JobName: fmt.Sprintf("job " + strconv.Itoa(i))}
		if !TryEnqueue(job, jobChan) {
			fmt.Println("max capacity reached, try later")
			close(jobChan)
			break
		} else {
			fmt.Printf("job %d assigned\n", i)
		}
	}
	time.Sleep(time.Second)

}
