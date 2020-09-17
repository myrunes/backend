// Package workerpool provides a simple implementation
// of goroutine "thread" pools.
package workerpool

import (
	"sync"
)

// Job defines a job function which will be executed
// in a worker getting passed the worker ID and
// parameters specified when pushing the job.
type Job func(workerId int, params ...interface{}) interface{}

// WorkerPool provides a simple "thread" pool
// implementation based on goroutines.
type WorkerPool struct {
	jobs    chan jobWrapper
	results chan interface{}
	wg      sync.WaitGroup
}

// jobWrapper wraps a Job and its specified
// parameters.
type jobWrapper struct {
	job    Job
	params []interface{}
}

// New creates a new instance of WorkerPool and
// spawns the defined number (size) of workers
// available waiting for jobs.
func New(size int) *WorkerPool {
	w := &WorkerPool{
		jobs:    make(chan jobWrapper),
		results: make(chan interface{}),
	}

	for i := 0; i < size; i++ {
		go w.spawnWorker(i)
	}

	return w
}

// Push enqueues a job with specified parameters which
// will be passed on executing the job.
func (w *WorkerPool) Push(job Job, params ...interface{}) {
	w.jobs <- jobWrapper{
		job:    job,
		params: params,
	}
}

// Close closes the Jobs channel so that the workers
// stop after executing all enqueued jobs.
// This is nessecary to be executed before WaitBlocking
// is called.
func (w *WorkerPool) Close() {
	close(w.jobs)
}

// Results returns the read-only results channel where
// executed job results are pushed in.
func (w *WorkerPool) Results() <-chan interface{} {
	return w.results
}

// WaitBlocking blocks until all jobs are finished.
func (w *WorkerPool) WaitBlocking() {
	w.wg.Wait()
	close(w.results)
}

// spawnWorker spawns a new worker with the passed
// worker id and starts listening for incomming jobs.
func (w *WorkerPool) spawnWorker(id int) {
	for job := range w.jobs {
		if job.job != nil {
			w.wg.Add(1)
			w.results <- job.job(id, job.params...)
			w.wg.Done()
		}
	}
}
