package queue

import (
	"sync"
)

// NewWorker creates a worker instance.
func NewWorker(jobs <-chan IJob) *Worker {
	worker := &Worker{
		jobs: jobs,
		quit: make(chan bool),
	}

	go worker.start()

	return worker
}

// Worker handles and executes jobs.
type Worker struct {
	jobs <-chan IJob
	quit chan bool
	once sync.Once
}

// Stop sends a signal to stopWorkersGracefully the worker.
func (w *Worker) Stop() <-chan bool {
	w.once.Do(func() {
		go func() {
			w.quit <- true
		}()
	})

	return w.quit
}

func (w *Worker) start() {
	for {
		select {
		case job := <-w.jobs:
			job.Do() // TODO add error handling anr retry policy

		case <-w.quit:
			close(w.quit)
			return
		}
	}
}
