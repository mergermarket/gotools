package tools

import "fmt"

//Worker Manages a number of concurrent workers
type Worker interface {
	Acquire() (release func())
	size() int
	available() int
}

//workerPool Manages a number of concurrent workers
type workerPool struct {
	throttle        chan int
	numberOfWorkers int
}

const badMaxConcurrentWorkersErrMsg = "the number of concurrent workers should be greater than %d"

//NewWorker Creates an instance of worker pool
func NewWorker(numberOfConcurrentWorkers uint8) (Worker, error) {
	numberOfWorkers := int(numberOfConcurrentWorkers)
	if numberOfWorkers == 0 {
		return nil, fmt.Errorf(badMaxConcurrentWorkersErrMsg, numberOfConcurrentWorkers)
	}

	throttle := make(chan int, numberOfWorkers)
	return &workerPool{throttle: throttle, numberOfWorkers: numberOfWorkers}, nil
}

//Acquire Acquires a single worker that can be released
func (w *workerPool) Acquire() (release func()) {
	w.throttle <- 1

	release = func() {
		<-w.throttle
	}

	return
}

//Size returns the number of workers
func (w *workerPool) size() int {
	return cap(w.throttle)
}

//Available returns the number of available workers
func (w *workerPool) available() int {
	return w.numberOfWorkers - len(w.throttle)
}
