package pool

import "fmt"

type Job struct {
	Payload interface{}
}

func (j Job) String() string {
	return fmt.Sprintf("%v", j.Payload)
}

type Result struct {
	Payload interface{}
	Error error
}

