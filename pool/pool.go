package pool

import (
	"log"
)

type Pool struct {
	NumWorkers int
	Jobs chan Job
	Results chan Result
	Stop chan bool
	Done chan bool
}

func NewPool(n int) *Pool {
	return &Pool{
		NumWorkers: n,
		Jobs: make(chan Job),
		Results: make(chan Result),
		Stop: make(chan bool),
		Done: make(chan bool)}
}

func (p *Pool) Submit(jobs []Job, function ExecFunc) {
	workers := make([]*Worker, 0)

	for i := 0; i < p.NumWorkers; i++ {
		worker := Worker{
			ID:i, Function:function,
			Done:p.Done, Stop:p.Stop, Pool:p}
		log.Printf("adding worker with ID=%d", worker.ID)
		workers = append(workers, &worker)
	}

	go func(){
		for _, job := range jobs {
			p.Jobs <- job
		}
		close(p.Jobs)
	}()

	for _, worker := range workers {
		log.Printf("submitting worker with ID=%d", worker.ID)
		go worker.Submit(p.Jobs, p.Results)
	}
}

func (p *Pool) Wait() []Result {
	results := make([]Result, 0)
	keepWorking := p.NumWorkers
	for {
		select {
		case result, ok := <- p.Results:
			if !ok {
				log.Printf("execution ended")
				return results
			} else {
				results = append(results, result)
			}
		case <- p.Done:
			keepWorking -= 1
			if keepWorking == 0 {
				close(p.Results)
			}
		}
	}
}
