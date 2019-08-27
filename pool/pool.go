package pool

import (
	"fmt"
	"log"
)

type Pool struct {
	ID int
	NumWorkers int
	Jobs chan Job
	Results chan Result
	Stop chan bool
	Done chan bool
}

var poolId = 0

func NewPool(n int) *Pool {
	p := &Pool{
		ID:poolId,
		NumWorkers: n,
		Jobs: make(chan Job),
		Results: make(chan Result),
		Stop: make(chan bool),
		Done: make(chan bool)}
	poolId += 1
	return p
}

func (p *Pool) Log(format string, args ...interface{}) {
	log.Printf("[pool:%d] %s", p.ID, fmt.Sprintf(format, args...))
}

func (p *Pool) Submit(jobs []Job, function ExecFunc) {
	workers := p.prepareWorkers(function)

	go func(){
		for _, job := range jobs {
			p.Jobs <- job
		}
		close(p.Jobs)
	}()

	p.spawnWorkers(workers)
}

func (p *Pool) Chain(prev <-chan Result, function ExecFunc) {
	workers := p.prepareWorkers(function)

	go func() {
		for result := range prev {
			if err := result.Error; err != nil {
				p.Log("previous result has error: %s", err.Error())
			} else {
				p.Jobs <- Job{Payload:result.Payload}
			}
		}
		close(p.Jobs)
	}()

	p.spawnWorkers(workers)
}

func (p *Pool) Wait() []Result {
	results := make([]Result, 0)
	keepWorking := p.NumWorkers
	for {
		select {
		case result, ok := <- p.Results:
			if !ok {
				p.Log("execution ended")
				return results
			} else {
				results = append(results, result)
			}
		case <- p.Done:
			keepWorking -= 1
			if keepWorking == 0 {
				p.Log("no more results expected, closing channel")
				close(p.Results)
			}
		}
	}
}

func (p *Pool) prepareWorkers(f ExecFunc) []*Worker {
	workers := make([]*Worker, 0)
	for i := 0; i < p.NumWorkers; i++ {
		worker := Worker{
			ID:i, Function:f,
			Done:p.Done, Stop:p.Stop, Pool:p}
		p.Log("adding worker with ID=%d", worker.ID)
		workers = append(workers, &worker)
	}
	return workers
}

func (p *Pool) spawnWorkers(workers []*Worker) {
	for _, worker := range workers {
		p.Log("submitting worker with ID=%d", worker.ID)
		go worker.Submit(p.Jobs, p.Results)
	}
}
