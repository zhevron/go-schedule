package schedule

import (
	"testing"
	"time"
)

func TestNewTrigger(t *testing.T) {
	trigger := NewTrigger()
	if trigger == nil {
		t.Error("Trigger was nil")
	}
}

func TestTrigger_Every(t *testing.T) {
	trigger := NewTrigger()
	trigger.Every("15m")
	if trigger.interval != 15*time.Minute {
		t.Errorf(
			"Interval does not match. Got %v, expected %v",
			trigger.interval.String(),
			(15 * time.Minute).String(),
		)
	}
}

func TestTrigger_EveryNegative(t *testing.T) {
	trigger := NewTrigger()
	trigger.Every("-15m")
	if trigger.interval != 0 {
		t.Errorf(
			"Interval does not match. Got %v, expected %v",
			trigger.interval.String(),
			0,
		)
	}
}

func TestTrigger_EveryPanic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Every did not panic with an invalid interval string")
		}
	}()
	trigger := NewTrigger()
	trigger.Every("15x")
}

func TestTrigger_Limit(t *testing.T) {
	trigger := NewTrigger()
	trigger.Limit(10)
	if trigger.limit != 10 {
		t.Errorf("Limit did not match. Got %v, expected 10", trigger.limit)
	}
}

func TestTrigger_LimitNegative(t *testing.T) {
	trigger := NewTrigger()
	trigger.Limit(-10)
	if trigger.limit != 0 {
		t.Errorf("Limit did not match. Got %v, expected 0", trigger.limit)
	}
}

func TestTrigger_Next(t *testing.T) {
	start := time.Now().Add(1 * time.Hour)
	next := start.Add(30 * time.Minute)
	trigger := &Trigger{
		interval: 30 * time.Minute,
		start:    start,
	}
	n := trigger.Next()
	if !n.Equal(next) {
		t.Errorf("Next time did not match. Got %v, expected %v", n, next)
	}
}

func TestTrigger_NextBefore(t *testing.T) {
	now := time.Now()
	start := now.AddDate(0, -1, 0)
	next := now.Add(30 * time.Minute)
	trigger := &Trigger{
		interval: 30 * time.Minute,
		start:    start,
	}
	n := trigger.Next()
	if !n.Equal(next) {
		t.Errorf("Next time did not match. Got %v, expected %v", n, next)
	}
}

func TestTrigger_NextLimit(t *testing.T) {
	now := time.Now()
	start := now.AddDate(0, -1, 0)
	next := start.Add(30 * time.Minute)
	trigger := &Trigger{
		interval: 30 * time.Minute,
		start:    start,
		limit:    1,
	}
	n := trigger.Next()
	if !n.Equal(next) {
		t.Errorf("Next time did not match. Got %v, expected %v.", n, next)
	}
}

func TestTrigger_NextZero(t *testing.T) {
	trigger := &Trigger{
		interval: 0,
		start:    time.Now(),
	}
	z := time.Time{}
	n := trigger.Next()
	if !n.IsZero() {
		t.Errorf("Next time did not match. Got %v, expected %v.", n, z)
	}
}

func TestTrigger_From(t *testing.T) {
	trigger := NewTrigger()
	tm := time.Now().Add(15 * time.Minute)
	trigger.From(tm)
	if !trigger.start.Equal(tm) {
		t.Errorf(
			"Start time does not match. Got %v, expected %v",
			trigger.start,
			tm,
		)
	}
}
