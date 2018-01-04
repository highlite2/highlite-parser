package queue

import (
	"sync"
)

// NewWorker creates a worker instance.
func NewWorker(jobs <-chan IJob) *Worker {
	worker := &Worker{
		jobs:    jobs,
		quit:    make(chan bool),
		stopped: make(chan bool),
	}

	go worker.start()

	return worker
}

// Worker handles and executes jobs.
type Worker struct {
	jobs    <-chan IJob
	quit    chan bool
	once    sync.Once
	stopped chan bool
}

// Stop sends a signal to stop the worker.
func (w *Worker) Stop() <-chan bool {
	w.once.Do(func() {
		close(w.quit)
	})

	return w.stopped
}

// Executes jobs.
func (w *Worker) start() {
	defer close(w.stopped)

	for {
		select {
		case job, ok := <-w.jobs:
			if !ok {
				return
			}

			job.Do() // TODO add error handling anr retry policy

		case <-w.quit:
			return
		}
	}
}
