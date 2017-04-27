package tools

import "testing"

func TestWorkerPool(t *testing.T) {

	t.Run("should Not create a worker if number of workers to create is less than one", func(t *testing.T) {
		numberOfWorkers := 0

		_, err := NewWorker(uint8(numberOfWorkers))

		if err == nil {
			t.Fatalf("should not have createe an instance of the worker %v", err)
		}
	})

	t.Run("should Acquire an worker and successfully release it back to the pool", func(t *testing.T) {
		numberOfWorkers := 5

		workerPool, err := NewWorker(uint8(numberOfWorkers))

		if err != nil {
			t.Fatalf("failed to create an instance of the worker %v", err)
		}

		releaseFirstWorker := workerPool.Acquire()

		if workerPool.size() != numberOfWorkers {
			t.Fatalf("expected %d but got %d", numberOfWorkers, workerPool.size())
		}

		availableWorkersBeforeRelease := numberOfWorkers - 1

		if workerPool.available() != availableWorkersBeforeRelease {
			t.Fatalf("expected %d but got %d", availableWorkersBeforeRelease, workerPool.available())
		}

		releaseFirstWorker()

		availableWorkersAfterRelease := numberOfWorkers

		if workerPool.available() != availableWorkersAfterRelease {
			t.Fatalf("expected %d but got %d", availableWorkersAfterRelease, workerPool.available())
		}

	})

}
