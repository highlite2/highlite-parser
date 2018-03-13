package queue

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorker_ExecutesJob(t *testing.T) {
	// arrange
	jobIsDone := make(chan bool)
	job := NewCallbackJob(func() error {
		jobIsDone <- true

		return nil
	})
	jobs := make(chan IJob)
	NewWorker(jobs)
	var err1, err2 error

	// act
	select {
	case jobs <- job:
	case <-time.After(time.Second):
		err1 = errors.New("can't push the job")
	}

	select {
	case <-jobIsDone:
	case <-time.After(time.Second):
		err2 = errors.New("job is not completed")
	}

	// assert
	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

func TestWorker_StopsWorker(t *testing.T) {
	// arrange
	jobs := make(chan IJob)
	worker := NewWorker(jobs)
	var err error

	// act
	select {
	case <-worker.Stop():
	case <-time.After(time.Second):
		err = errors.New("can't stop the worker")
	}

	// assert
	assert.Nil(t, err)
}

func TestWorker_JobIsFinishedAfterStop(t *testing.T) {
	// arrange
	jobs := make(chan IJob)
	worker := NewWorker(jobs)
	var jobIsProcessed bool
	job := NewCallbackJob(func() error {
		jobIsProcessed = true
		return nil
	})

	// act
	jobs <- job
	<-worker.Stop()

	// assert
	assert.True(t, jobIsProcessed)

}
