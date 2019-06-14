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
