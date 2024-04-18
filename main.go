
package main

import (
	"fmt"
	"runtime"
	"time"
)

const FIB_COUNT = 45 // the number of fibonacci numbers to calculate
const WORKER_COUNT = 100

func main() {
	// create channels for our jobs and results
	jobs := make(chan int, FIB_COUNT)
	results := make(chan int, FIB_COUNT)
	done := make(chan bool)

	job_counter := FIB_COUNT

	// create a worker pool
	for range WORKER_COUNT {
		go worker(jobs, results, done, &job_counter)
	}

	fmt.Printf("Total number of Go Routine: %d\n", runtime.NumGoroutine())

	// populate our queue with jobs and then close the channel
	for n := range FIB_COUNT {
		jobs <- n
	}
	fmt.Printf("%d jobs enqueued\n\n", len(jobs))
	close(jobs)

	start := time.Now() // start timer

	<-done // listen for the done signal
	close(results)

	elapsed := time.Since(start) // end timer
	fmt.Printf("Processing took %s\n", elapsed)

}

// worker function that will calculate the fibonacci number
func worker(jobs <-chan int, results chan<- int, done chan<- bool, job_counter *int) {
	for n := range jobs {
		f := fib(n)
		fmt.Println(f)
		*job_counter--
		results <- f
	}
	if *job_counter == 0 {
		done <- true // signal that we're done
	}
}

// calculate the fibonacci number
func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
