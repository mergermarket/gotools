package tools

import "fmt"

//Worker Manages a number of concurrent workers
type Worker interface {
	Acquire() (release func())
}

//workerPool Manages a number of concurrent workers
type workerPool struct {
	throttle chan int
}

const badMaxConcurrentWorkersErrMsg = "the number of concurrent workers should be greater than %d"

//NewWorker Creates an instance of worker pool
func NewWorker(numberOfConcurrentWorkers uint8) (Worker, error) {
	if int(numberOfConcurrentWorkers) == 0 {
		return nil, fmt.Errorf(badMaxConcurrentWorkersErrMsg, numberOfConcurrentWorkers)
	}

	throttle := make(chan int, int(numberOfConcurrentWorkers))
	return &workerPool{throttle: throttle}, nil
}

//Acquire Acquires a single worker that can be released
func (w *workerPool) Acquire() (release func()) {
	w.throttle <- 1

	release = func() {
		<-w.throttle
	}

	return
}
