package schedule

import (
	"sync"
	"time"
)

// A Queue represents a Job queue, responsible for running Jobs when scheduled.
type Queue struct {
	Jobs      []*Job
	errors    chan JobError
	mutex     sync.RWMutex
	results   chan JobResult
	suspended bool
}

// NewQueue creates a new Queue.
// By default, the Queue is initialized with a max results and error
// buffer of 10. See MaxBufferedErrors and MaxBufferedResults.
func NewQueue() *Queue {
	return &Queue{
		Jobs:      make([]*Job, 0),
		errors:    make(chan JobError, 10),
		results:   make(chan JobResult, 10),
		suspended: false,
	}
}

// Add appends a job to this queue.
// If the job is already present, the function returns without adding it.
func (q *Queue) Add(job *Job) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for _, j := range q.Jobs {
		if job == j {
			return
		}
	}
	q.Jobs = append(q.Jobs, job)
}

// Errors returns the channel on which job errors are emitted.
func (q *Queue) Errors() chan JobError {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return q.errors
}

// MaxBufferedErrors sets the buffer length of the job errors channel.
func (q *Queue) MaxBufferedErrors(n int) {
	errors := make(chan JobError, n)
	q.mutex.Lock()
	defer q.mutex.Unlock()
	close(q.errors)
	for err := range q.errors {
		select {
		case errors <- err:
		}
	}
	q.errors = errors
}

// MaxBufferedResults sets the buffer length of the job results channel.
func (q *Queue) MaxBufferedResults(n int) {
	results := make(chan JobResult, n)
	close(q.results)
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for result := range q.results {
		select {
		case results <- result:
		}
	}
	q.results = results
}

// Results returns the channel on which job results are emitted.
func (q *Queue) Results() chan JobResult {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return q.results
}

// Resume will resume the suspended queue and jobs will be checked next Run.
func (q *Queue) Resume() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.suspended = false
}

// Run checks all jobs if they should be run and triggers each of them in
// their own goroutine. Job results and errors are emitted to the Results and
// Errors channels respectively.
// If the queue is suspended, no jobs are checked.
func (q *Queue) Run() {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	if !q.suspended {
		for _, job := range q.Jobs {
			next := job.NextRun()
			if !next.IsZero() && !next.After(time.Now()) {
				go func(job *Job) {
					res, err := job.Run()
					if len(res) > 0 {
						q.mutex.Lock()
						select {
						case q.results <- JobResult{job.Name, res}:
						}
						q.mutex.Unlock()
					}
					if err != nil {
						q.mutex.Lock()
						select {
						case q.errors <- JobError{job.Name, err}:
						}
						q.mutex.Unlock()
					}
				}(job)
			}
		}
	}
}

// Suspend will suspend the queue and no jobs will be run until Resumed.
func (q *Queue) Suspend() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.suspended = true
}

// Suspended returns whether the queue is suspended.
func (q *Queue) Suspended() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return q.suspended
}
