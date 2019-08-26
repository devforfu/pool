package main

import (
	"log"
	"math/rand"
	"pool/pool"
	"strings"
)

func main() {
	jobs := makeJobs(10)

	lengthCounter := pool.NewPool(2)
	lengthCounter.Submit(
		jobs,
		func(job pool.Job, _ *pool.Worker) pool.Result {
			return pool.Result{Payload: len(job.Payload.(string))}
		},
	)

	isEvenLengthChecker := pool.NewPool(2)
	isEvenLengthChecker.Chain(
		lengthCounter.Results,
		func(job pool.Job, _ *pool.Worker) pool.Result {
			return pool.Result{Payload: job.Payload.(int) % 2 == 0}
		},
	)

	results := isEvenLengthChecker.Wait()
	log.Printf("%v", results)

	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	for result := range isEvenLengthChecker.Results {
	//		log.Printf("Is Even? %t", result.Payload.(bool))
	//	}
	//	wg.Done()
	//}()
	//
	//wg.Wait()

	//results := p.Wait()
	//lengths := make([]int, 0)
	//for _, result := range results {
	//	lengths = append(lengths, result.Payload.(int))
	//}
	//log.Printf("collected results: %v", lengths)
	//log.Printf("completed")
}

func randomString(n int) string {
	const domain = "abcdefghijklmnop"
	var b strings.Builder
	for i := 0; i < n; i++ {
		index := rand.Intn(len(domain))
		b.WriteString(string(domain[index]))
	}
	return b.String()
}

func makeJobs(n int) []pool.Job {
	jobs := make([]pool.Job, 0)
	for i := 0; i < 100; i++ {
		size := rand.Intn(10) + 5
		jobs = append(jobs, pool.Job{Payload:randomString(size)})
	}
	return jobs
}
