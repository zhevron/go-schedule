// Package schedule implements task scheduling via goroutines.
package schedule

import "time"

// A Trigger represents the time schedule for a job.
type Trigger struct {
	interval time.Duration
	start    time.Time
	limit    int64
}

// NewTrigger creates a new Trigger.
// By default, the From field is populated with the current time.
func NewTrigger() *Trigger {
	return &Trigger{
		interval: 0,
		start:    time.Now(),
		limit:    0,
	}
}

// Every sets the recurrence for the Trigger. The value passed needs to be
// parsable by time.ParseDuration. E.g Trigger.Every("1m30s").
// If the provided string cannot be parsed, the function will panic with
// the error.
// If a negative interval is provided, the interval will be set to 0 and the
// trigger will always return a zeroed time.Time from Next.
func (t *Trigger) Every(str string) *Trigger {
	d, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}
	if d < 0 {
		d = 0
	}
	t.interval = d
	return t
}

// Limit sets the number of times a job is allowed to run before Next returns
// a zeroed time.Time.
// If 0 or a negative value is provided, the limit will be set to 0 (no limit).
func (t *Trigger) Limit(n int64) *Trigger {
	if n < 0 {
		n = 0
	}
	t.limit = n
	return t
}

// Next returns the next scheduled time for the Trigger, counting from the
// current time. If the interval is 0, a zeroed time.Time will be returned.
func (t *Trigger) Next() time.Time {
	if t.interval == 0 {
		return time.Time{}
	}
	now := time.Now()
	next := t.start.Add(t.interval)
	current := int64(0)
	for next.Before(now) {
		current = current + 1
		if t.limit > 0 && current >= t.limit {
			return next
		}
		next = next.Add(t.interval)
	}
	return next
}

// From sets the start time from which the recurrence is counted from.
func (t *Trigger) From(tm time.Time) *Trigger {
	t.start = tm
	return t
}
