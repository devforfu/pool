package main

import (
	"log"
	"math/rand"
	"pool/pool"
	"strings"
)

func randomString(n int) string {
	const domain = "abcdefghijklmnop"
	var b strings.Builder
	for i := 0; i < n; i++ {
		index := rand.Intn(len(domain))
		b.WriteString(string(domain[index]))
	}
	return b.String()
}

func main() {
	jobs := make([]pool.Job, 0)
	for i := 0; i < 100; i++ {
		size := rand.Intn(10) + 5
		if rand.Intn(2) == 0 {
			jobs = append(jobs, pool.Job{Payload:10})
		} else {
			jobs = append(jobs, pool.Job{Payload:randomString(size)})
		}
	}
	p := pool.NewPool(3)
	p.Submit(jobs, func(job pool.Job, _ *pool.Worker) pool.Result {
		return pool.Result{Payload: len(job.Payload.(string))}
	})
	results := p.Wait()
	lengths := make([]int, 0)
	for _, result := range results {
		lengths = append(lengths, result.Payload.(int))
	}
	log.Printf("collected results: %v", lengths)
	log.Printf("completed")
}
