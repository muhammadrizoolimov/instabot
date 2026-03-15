package worker

import (
	"log"
	"sync"
)

const maxWorkers = 5

type Job struct {
	ID      string
	URL     string
	Handler func(string) error
}

type WorkerPool struct {
	jobs    chan Job
	wg      sync.WaitGroup
	mu      sync.Mutex
	active  map[string]bool
}

var (
	pool     *WorkerPool
	poolOnce sync.Once
)

// GetPool returns the singleton worker pool instance
func GetPool() *WorkerPool {
	poolOnce.Do(func() {
		pool = &WorkerPool{
			jobs:   make(chan Job, 100),
			active: make(map[string]bool),
		}
		pool.start()
	})
	return pool
}

// start initializes worker goroutines
func (wp *WorkerPool) start() {
	for i := 0; i < maxWorkers; i++ {
		go wp.worker(i)
	}
}

// worker processes jobs from the queue
func (wp *WorkerPool) worker(id int) {
	for job := range wp.jobs {
		log.Printf("Worker %d processing job: %s", id, job.ID)

		wp.mu.Lock()
		wp.active[job.ID] = true
		wp.mu.Unlock()

		if err := job.Handler(job.URL); err != nil {
			log.Printf("Worker %d job %s failed: %v", id, job.ID, err)
		}

		wp.mu.Lock()
		delete(wp.active, job.ID)
		wp.mu.Unlock()

		wp.wg.Done()
	}
}

// Submit adds a job to the worker pool queue
func (wp *WorkerPool) Submit(job Job) {
	wp.mu.Lock()
	if wp.active[job.ID] {
		wp.mu.Unlock()
		log.Printf("Job %s already in queue, skipping", job.ID)
		return
	}
	wp.mu.Unlock()

	wp.wg.Add(1)
	wp.jobs <- job
}

// Wait waits for all jobs to complete
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Close closes the worker pool
func (wp *WorkerPool) Close() {
	close(wp.jobs)
}
