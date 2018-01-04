package queue

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
)

func TestPool_AllJobsAreExecuted(t *testing.T) {
	// arrange
	jobNumber := 10000
	pool := NewPool(10)
	ready := make(chan bool)
	job := NewCallbackJob(func() error {
		ready <- true
		return nil
	})
	counter := 0

	// act
	for i := 0; i < jobNumber; i++ {
		pool.AddJob(job)
	}
	for i := 0; i < jobNumber; i++ {
		select {
		case <-ready:
			counter++
		case <-time.After(time.Second):
			i = jobNumber
			counter = -1
		}
	}

	// assert
	assert.Equal(t, jobNumber, counter)
}

func TestPool_AddedJobsAreExecutedAfterStop(t *testing.T) {
	// arrange
	jobNumber := 10000
	var jobDone, jobAdded int32
	var err error
	pool := NewPool(10)
	addResult := make(chan bool, jobNumber)

	// act
	for i := 0; i < jobNumber; i++ {
		go func(t int) {
			addResult <- <-pool.AddJob(NewCallbackJob(func() error {
				atomic.AddInt32(&jobDone, 1)
				return nil
			}))
		}(i)
	}

	stopped := pool.Stop()

	for i := 0; i < jobNumber; i++ {
		if <-addResult {
			jobAdded++
		}
	}

	select {
	case <-stopped:
	case <-time.After(time.Second):
		err = errors.New("can't stop the pool")
	}

	// assert
	assert.Nil(t, err)
	assert.Equal(t, jobAdded, jobDone)
}
