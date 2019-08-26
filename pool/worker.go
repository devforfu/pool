package pool

import (
	"fmt"
	"log"
)

type Worker struct {
	ID int
	Function ExecFunc
	Pool *Pool
	Done chan<- bool
	Stop <-chan bool
}

type ExecFunc func(Job, *Worker) Result

func (w *Worker) Submit(jobs <-chan Job, results chan<- Result) {
	defer func() {
		w.Done <- true
		if err := recover(); err != nil {
			log.Printf("*** worker %d paniced! ***", w.ID)
			log.Printf("%v", err)
		}
	}()
	running := true
	for running {
		select {
		case job, ok := <- jobs:
			if !ok {
				w.Log("setting the channel to nil")
				jobs = nil
				running = false
			} else {
				result := w.Function(job, w)
				w.Log("job %v -> result %v", job, result)
				results <- result
			}
		case <- w.Stop:
			w.Log("interrupted!")
			running = false
		}
	}
}

func (w *Worker) Log(format string, args ...interface{}) {
	log.Printf("[worker:%d] %s", w.ID, fmt.Sprintf(format, args...))
}