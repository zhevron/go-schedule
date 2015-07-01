package schedule

import (
	"fmt"
	"testing"
	"time"
)

func TestNewScheduler(t *testing.T) {
	s := NewScheduler()
	if s == nil {
		t.Error("Scheduler was nil")
	}
}

func TestScheduler_Add(t *testing.T) {
	s := NewScheduler()
	j, err := NewJob("test", func() { return })
	if err != nil {
		t.Fatalf("Could not create test Job: %v", err)
	}
	s.Add(j)
	if len(s.Queues["default"].Jobs) != 1 {
		t.Errorf(
			"Number of jobs in default queue did not match. Got %d, expected 1",
			len(s.Queues["default"].Jobs),
		)
	}
}

func TestScheduler_AddToQueue(t *testing.T) {
	s := NewScheduler()
	j, err := NewJob("test", func() { return })
	if err != nil {
		t.Fatalf("Could not create test Job: %v", err)
	}
	s.AddToQueue("default", j)
	if len(s.Queues["default"].Jobs) != 1 {
		t.Errorf(
			"Number of jobs in default queue did not match. Got %d, expected 1",
			len(s.Queues["default"].Jobs),
		)
	}
}

func TestScheduler_Errors(t *testing.T) {
	s := NewScheduler()
	if s.Errors() == nil {
		t.Error("Error buffer was nil.")
	}
}

func TestScheduler_MaxBufferedErrors(t *testing.T) {
	s := NewScheduler()
	s.MaxBufferedErrors(50)
	if cap(s.errors) != 50 {
		t.Errorf(
			"Error buffer for Scheduler did not match. Got %d, expected 50.",
			cap(s.errors),
		)
	}
	if cap(s.Queues["default"].errors) != 50 {
		t.Errorf(
			"Error buffer for default queue did not match. Got %d, expected 50.",
			cap(s.Queues["default"].errors),
		)
	}
}

func TestScheduler_MaxBufferedErrorsCopy(t *testing.T) {
	s := NewScheduler()
	s.errors <- JobError{"test", nil}
	s.MaxBufferedErrors(50)
	if len(s.errors) != 1 {
		t.Errorf(
			"Number of errors in buffer did not match. Got %d, expected 1.",
			len(s.errors),
		)
	}
}

func TestScheduler_MaxBufferedResults(t *testing.T) {
	s := NewScheduler()
	s.MaxBufferedResults(50)
	if cap(s.results) != 50 {
		t.Errorf(
			"Result buffer for Scheduler did not match. Got %d, expected 50.",
			cap(s.results),
		)
	}
	if cap(s.Queues["default"].results) != 50 {
		t.Errorf(
			"Result buffer for default queue did not match. Got %d, expected 50.",
			cap(s.Queues["default"].results),
		)
	}
}

func TestScheduler_MaxBufferedResultsCopy(t *testing.T) {
	s := NewScheduler()
	s.results <- JobResult{"test", []interface{}{}}
	s.MaxBufferedResults(50)
	if len(s.results) != 1 {
		t.Errorf(
			"Number of results in buffer did not match. Got %d, expected 1.",
			len(s.results),
		)
	}
}

func TestScheduler_Results(t *testing.T) {
	s := NewScheduler()
	if s.Results() == nil {
		t.Error("Result buffer was nil.")
	}
}

func TestScheduler_Start(t *testing.T) {
	s := NewScheduler()
	defer s.Stop()
	if err := s.Start(); err != nil {
		t.Errorf("Scheduler errored on Start: %v", err)
	}
	if !s.Running() {
		t.Error("Scheduler is not running after Start")
	}
	j1, err := NewJob("test1", func() string { return "test" })
	if err != nil {
		t.Fatalf("Could not create test Job: %v", err)
	}
	j1.Schedule().Every("1ms").Limit(1)
	s.Add(j1)
	j2, err := NewJob("test2", func() { panic("test") })
	if err != nil {
		t.Fatalf("Could not create test Job: %v", err)
	}
	j2.Schedule().Every("1ms").Limit(1)
	s.Add(j2)
	time.Sleep(100 * time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	if len(s.errors) != 1 {
		fmt.Println(s.errors)
		t.Errorf(
			"Errors in buffer did not match. Got %d, expected 1",
			len(s.errors),
		)
	}
	if len(s.results) != 1 {
		fmt.Println(s.results)
		t.Errorf(
			"Results in buffer did not match. Got %d, expected 1",
			len(s.results),
		)
	}
}

func TestScheduler_StartRunning(t *testing.T) {
	s := NewScheduler()
	defer s.Stop()
	if err := s.Start(); err != nil {
		t.Errorf("Scheduler errored on Start: %v", err)
	}
	if err := s.Start(); err == nil {
		t.Error("Scheduler did not error on Start while already running")
	}
}

func TestScheduler_Stop(t *testing.T) {
	s := NewScheduler()
	if err := s.Start(); err != nil {
		t.Errorf("Scheduler errored on Start: %v", err)
	}
	if !s.Running() {
		t.Error("Scheduler is not running after Start")
	}
	s.Stop()
	if s.Running() {
		t.Error("Scheduler is not stopped after Stop")
	}
}

func TestScheduler_Queue(t *testing.T) {
	s := NewScheduler()
	s.Queue("test", NewQueue())
	if len(s.Queues) != 2 {
		t.Errorf(
			"Number of queues in Scheduler did not match. Got %d, expected 2.",
			len(s.Queues),
		)
	}
	if _, ok := s.Queues["test"]; !ok {
		t.Error("Could not find test queue in Scheduler.")
	}
}
