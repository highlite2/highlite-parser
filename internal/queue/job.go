package queue

// IJob is a common job interface
type IJob interface {
	Do() error
}
