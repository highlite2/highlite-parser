package queue

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestWorker_ExecutesJob(t *testing.T) {
	// arrange
	jobIsDone := make(chan bool)
	jobs := make(chan IJob)
	job := NewCallbackJob(func() error {
		jobIsDone <- true

		return nil
	})
	NewWorker(jobs)
	var err error

	// act
	jobs <- job
	select {
	case <-jobIsDone:
	default:
		err = errors.New("job is not completed")
	}

	// assert
	assert.Nil(t, err)
}
