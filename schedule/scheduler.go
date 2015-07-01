// Package schedule implements task scheduling via goroutines.
package schedule

import (
	"errors"
	"sync"
	"time"
)

// A Scheduler represents an active Queue runner.
type Scheduler struct {
	Queues  map[string]*Queue
	errors  chan JobError
	mutex   sync.RWMutex
	results chan JobResult
	running bool
}

// NewScheduler creates a new Scheduler with a single "default" queue.
// By default, the Scheduler is initialized with a max results and error
// buffer of 10. See MaxBufferedErrors and MaxBufferedResults.
func NewScheduler() *Scheduler {
	return &Scheduler{
		Queues: map[string]*Queue{
			"default": NewQueue(),
		},
		errors:  make(chan JobError, 10),
		results: make(chan JobResult, 10),
		running: false,
	}
}

// Add appends a job to the "default" queue of this Scheduler.
func (s *Scheduler) Add(job *Job) {
	s.AddToQueue("default", job)
}

// AddToQueue appends a job to the given queue of this Scheduler.
func (s *Scheduler) AddToQueue(queue string, job *Job) {
	s.mutex.RLock()
	s.Queues[queue].Add(job)
	s.mutex.RUnlock()
}

// Errors returns the channel on which job errors are emitted.
func (s *Scheduler) Errors() chan JobError {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.errors
}

// MaxBufferedErrors sets the buffer length of the job errors channel.
// This method will also be called on all available queues.
// Please note that this will not affect any channels added after.
func (s *Scheduler) MaxBufferedErrors(n int) {
	errors := make(chan JobError, n)
	s.mutex.Lock()
	close(s.errors)
	for err := range s.errors {
		select {
		case errors <- err:
		}
	}
	s.errors = errors
	s.mutex.Unlock()
	s.mutex.RLock()
	for _, queue := range s.Queues {
		queue.MaxBufferedErrors(n)
	}
	s.mutex.RUnlock()
}

// MaxBufferedResults sets the buffer length of the job results channel.
// This method will also be called on all available queues.
// Please note that this will not affect any channels added after.
func (s *Scheduler) MaxBufferedResults(n int) {
	results := make(chan JobResult, n)
	s.mutex.Lock()
	close(s.results)
	for result := range s.results {
		select {
		case results <- result:
		}
	}
	s.results = results
	s.mutex.Unlock()
	s.mutex.RLock()
	for _, queue := range s.Queues {
		queue.MaxBufferedResults(n)
	}
	s.mutex.RUnlock()
}

// Results returns the channel on which job results are emitted.
func (s *Scheduler) Results() chan JobResult {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.results
}

// Running returns whether this Scheduler is currently running.
func (s *Scheduler) Running() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.running
}

// Start begins the process of running the queues.
// If this Scheduler has been started, an error is returned.
// All job results in the Queues proxied into the Scheduler channels.
func (s *Scheduler) Start() error {
	if s.Running() {
		return errors.New("schedule: scheduler is already running")
	}
	s.mutex.Lock()
	s.running = true
	s.mutex.Unlock()
	go func() {
		for s.running {
			for _, queue := range s.Queues {
				go func(queue *Queue) {
					queue.Run()
				}(queue)
				go func(queue *Queue) {
					for len(queue.Errors()) > 0 || len(queue.Results()) > 0 {
						select {
						case err := <-queue.Errors():
							s.mutex.Lock()
							select {
							case s.errors <- err:
							}
							s.mutex.Unlock()
						case res := <-queue.Results():
							s.mutex.Lock()
							select {
							case s.results <- res:
							}
							s.mutex.Unlock()
						}
					}
				}(queue)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return nil
}

// Stop tells the Scheduler to stop processing queues after the current run.
// If the Scheduler is not running, this will have no effect.
func (s *Scheduler) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.running = false
}

// Queue adds a new Queue to this Scheduler.
// If a Queue is already present with the same name, it will be overwritten.
// Please note that any calls to MaxBufferedErrors or MaxBufferedResults does
// not affect any Queues added later. You will need to call these manually.
func (s *Scheduler) Queue(name string, queue *Queue) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Queues[name] = queue
}
