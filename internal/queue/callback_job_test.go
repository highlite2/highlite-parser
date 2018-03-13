package queue

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallbackJob_ExecuteJob_NoError(t *testing.T) {
	// arrange
	callback := func() error {
		return nil
	}
	var job IJob = NewCallbackJob(callback)

	// act
	err := job.Do()

	// assert
	assert.Nil(t, err)
}

func TestCallbackJob_ExecuteJob_WithError(t *testing.T) {
	// arrange
	expectedError := errors.New("some error")
	callback := func() error {
		return expectedError
	}
	var job IJob = NewCallbackJob(callback)

	// act
	err := job.Do()

	// assert
	assert.ObjectsAreEqual(expectedError, err)
}
