package queue

import "fmt"

var _ IJob = (*CallbackJob)(nil)

// NewCallbackJob creates new CallbackJob instance.
func NewCallbackJob(callback func() error) *CallbackJob {
	return &CallbackJob{
		callback: callback,
	}
}

// CallbackJob is a job, that is meant for executing given callback.
type CallbackJob struct {
	callback func() error
}

// Do executes the job
func (c *CallbackJob) Do() error {
	if c.callback == nil {
		return fmt.Errorf("callback is nil")
	}

	return c.callback()
}
