package schedule

import (
	"errors"
	"strings"
	"testing"
)

func TestNewJob(t *testing.T) {
	j, err := NewJob("test", func() { return })
	if err != nil {
		t.Fatal(err)
	}
	if j == nil {
		t.Error("Job was nil")
	}
}

func TestNewJobError(t *testing.T) {
	_, err := NewJob("test", 1)
	if err == nil || !strings.Contains(err.Error(), "function") {
		t.Errorf("Job was created with non-function parameter. Error: %#q", err)
	}
}

func TestJob_Args(t *testing.T) {
	a := []interface{}{1, "asd", 2, "def"}
	j, err := NewJob("test", func() { return }, a...)
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Args()) != len(a) {
		t.Fatalf("Argument count did not match. Expected %d, got %d", len(a), len(j.Args()))
	}
	for i, arg := range j.Args() {
		if arg != a[i] {
			t.Errorf("Argument %d did not match. Expected %#q, got %#q", i, a[i], arg)
		}
	}
}

func TestJob_LastRun(t *testing.T) {
	j, err := NewJob("test", func() { return })
	if err != nil {
		t.Fatal(err)
	}
	j.Run()
	if j.LastRun().IsZero() {
		t.Error("NextRun returned a zeroed time value")
	}
}

func TestJob_NextRun(t *testing.T) {
	j, err := NewJob("test", func() { return })
	if err != nil {
		t.Fatal(err)
	}
	j.Schedule().Every("15m")
	if j.NextRun().IsZero() {
		t.Error("NextRun returned a zeroed time value")
	}
}

func TestJob_Run(t *testing.T) {
	runCount := 0
	defer func() {
		if runCount <= 0 {
			t.Errorf("Job ran a total of %d times", runCount)
		}
	}()
	j, err := NewJob("test", func() {
		runCount++
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = j.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestJob_RunArgs(t *testing.T) {
	j, err := NewJob("test", func(a int, b int) {
		return
	}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = j.Run(); err == nil || !strings.Contains(err.Error(), "arguments") {
		t.Errorf("Job did not fail with invalid arguments. Error: %#q", err)
	}
}

func TestJob_RunPanicError(t *testing.T) {
	j, err := NewJob("test", func() {
		panic(errors.New("test"))
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = j.Run(); err == nil || !strings.Contains(err.Error(), "panic") {
		t.Errorf("Job did not panic. Error: %#q", err)
	}
}

func TestJob_RunPanicValue(t *testing.T) {
	j, err := NewJob("test", func() {
		panic("test")
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = j.Run(); err == nil || !strings.Contains(err.Error(), "panic") {
		t.Errorf("Job did not panic. Error: %#q", err)
	}
}

func TestJob_Schedule(t *testing.T) {
	j, err := NewJob("test", func() { return })
	if err != nil {
		t.Fatal(err)
	}
	s := j.Schedule()
	if s == nil {
		t.Errorf("Schedule returned a nil Trigger")
	}
}
