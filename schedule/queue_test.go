package schedule

import (
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {
	q := NewQueue()
	if q == nil {
		t.Error("Queue was nil")
	}
}

func TestQueue_Add(t *testing.T) {
	q := NewQueue()
	j, err := NewJob("test", func() { return })
	if err != nil {
		t.Fatalf("Could not create test Job: %v", err)
	}
	q.Add(j)
	if len(q.Jobs) != 1 {
		t.Errorf(
			"Number of jobs in queue did not match. Got %d, expected 1",
			len(q.Jobs),
		)
	}
}

func TestQueue_AddDuplicate(t *testing.T) {
	q := NewQueue()
	j, err := NewJob("test", func() { return })
	if err != nil {
		t.Fatalf("Could not create test Job: %v", err)
	}
	q.Add(j)
	q.Add(j)
	if len(q.Jobs) != 1 {
		t.Errorf(
			"Number of jobs in queue did not match. Got %d, expected 1",
			len(q.Jobs),
		)
	}
}

func TestQueue_Errors(t *testing.T) {
	q := NewQueue()
	if q.Errors() == nil {
		t.Error("Error buffer was nil.")
	}
}

func TestQueue_MaxBufferedErrors(t *testing.T) {
	q := NewQueue()
	q.MaxBufferedErrors(50)
	if cap(q.errors) != 50 {
		t.Errorf(
			"Error buffer for Queue did not match. Got %d, expected 50.",
			cap(q.errors),
		)
	}
}

func TestQueue_MaxBufferedErrorsCopy(t *testing.T) {
	q := NewQueue()
	q.errors <- JobError{"test", nil}
	q.MaxBufferedErrors(50)
	if len(q.errors) != 1 {
		t.Errorf(
			"Number of errors in buffer did not match. Got %d, expected 1.",
			len(q.errors),
		)
	}
}

func TestQueue_MaxBufferedResults(t *testing.T) {
	q := NewQueue()
	q.MaxBufferedResults(50)
	if cap(q.results) != 50 {
		t.Errorf(
			"Result buffer for Queue did not match. Got %d, expected 50.",
			cap(q.results),
		)
	}
}

func TestQueue_MaxBufferedResultsCopy(t *testing.T) {
	q := NewQueue()
	q.results <- JobResult{"test", []interface{}{}}
	q.MaxBufferedResults(50)
	if len(q.results) != 1 {
		t.Errorf(
			"Number of results in buffer did not match. Got %d, expected 1.",
			len(q.results),
		)
	}
}

func TestQueue_Results(t *testing.T) {
	q := NewQueue()
	if q.Results() == nil {
		t.Error("Result buffer was nil.")
	}
}

func TestQueue_Resume(t *testing.T) {
	q := NewQueue()
	q.Resume()
	if q.Suspended() {
		t.Error("Queue is suspended after Resume")
	}
}

func TestQueue_Run(t *testing.T) {
	q := NewQueue()
	j1, err := NewJob("test", func() string { return "test" })
	if err != nil {
		t.Fatalf("Could not create test Job: %v", err)
	}
	j1.Schedule().Every("1ms").Limit(1)
	q.Add(j1)
	j2, err := NewJob("test", func() { panic("test") })
	if err != nil {
		t.Fatalf("Could not create test Job: %v", err)
	}
	j2.Schedule().Every("1ms").Limit(1)
	q.Add(j2)
	time.Sleep(100 * time.Millisecond)
	q.Run()
	time.Sleep(100 * time.Millisecond)
	if len(q.errors) != 1 {
		t.Errorf(
			"Errors in buffer did not match. Got %d, expected 1",
			len(q.errors),
		)
	}
	if len(q.results) != 1 {
		t.Errorf(
			"Results in buffer did not match. Got %d, expected 1",
			len(q.results),
		)
	}
}

func TestQueue_Suspend(t *testing.T) {
	q := NewQueue()
	q.Suspend()
	if !q.Suspended() {
		t.Error("Queue is active after Suspend")
	}
}
