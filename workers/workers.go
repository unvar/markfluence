package workers

import (
	"sync"
)

// LoadJobs is a util to loop thorugh the file list
// and put it into a channel
func LoadJobs(files []string, jobs chan<- string) {
	for _, file := range files {
		jobs <- file
	}
	close(jobs)
}

// CreateWorkerPool can be used to create a set of
// workers that run the submitted job function for all
// the jobs in a given channel
func CreateWorkerPool(
	size int,
	jobs <-chan string,
	handler func(payload string),
	done chan<- bool,
) {
	var wg sync.WaitGroup
	for i := 0; i < size; i++ {
		wg.Add(1)
		go worker(&wg, jobs, handler)
	}
	wg.Wait()

	done <- true
}

func worker(wg *sync.WaitGroup, jobs <-chan string, handler func(payload string)) {
	defer wg.Done()
	for job := range jobs {
		handler(job)
	}
}
