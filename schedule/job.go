package schedule

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// A JobResult represents the return value(s) from a job function.
// The name of the job is stored in .Name and the return values are stored
// as an array of interface{} in .Results
type JobResult struct {
	Name    string
	Results []interface{}
}

// A JobError represents a job execution error.
// THe name of the job is stored in .Name and the error in .Error
type JobError struct {
	Name  string
	Error error
}

// A Job represents an executable job.
type Job struct {
	Name     string
	function reflect.Value
	args     []reflect.Value
	trigger  *Trigger
	last     time.Time
}

// NewJob creates a new Job for the given function.
// Each argument passed after the function is passed along as an argument to
// the function when called.
func NewJob(name string, fun interface{}, args ...interface{}) (*Job, error) {
	function := reflect.ValueOf(fun)
	if function.Kind() != reflect.Func {
		return nil, errors.New("schedule: jobs can only be created for functions")
	}
	arguments := make([]reflect.Value, len(args))
	for i, arg := range args {
		arguments[i] = reflect.ValueOf(arg)
	}
	return &Job{
		Name:     name,
		function: function,
		args:     arguments,
		trigger:  NewTrigger(),
	}, nil
}

// Args returns the arguments associated with the job function.
func (j *Job) Args() []interface{} {
	args := make([]interface{}, len(j.args))
	for i, arg := range j.args {
		args[i] = arg.Interface()
	}
	return args
}

// LastRun returns the timestamp of the last successful run.
// Note that if the job errors, this timestamp will not be updated.
func (j *Job) LastRun() time.Time {
	return j.last
}

// NextRun returns the timestamp of the next scheduled run.
func (j *Job) NextRun() time.Time {
	next := j.trigger.Next()
	if j.last.Before(next) {
		return next
	}
	return time.Time{}
}

// Run attempts to call the job function with the provided arguments.
// If the arguments do not match the function, an error is returned.
// This function will also recover from any panic caused inside a job and
// return the panic value as an error.
func (j *Job) Run() (result []interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			switch e.(type) {
			case error:
				err = fmt.Errorf("schedule: job panicked with error %#q", e.(error))
			default:
				err = fmt.Errorf("schedule: job panicked with value %#q", e)
			}
		} else {
			j.last = time.Now()
		}
	}()
	for _, res := range j.function.Call(j.args) {
		result = append(result, res.Interface())
	}
	return
}

// Schedule creates a new Trigger and returns it so that a schedule may
// be constructed.
func (j *Job) Schedule() *Trigger {
	j.trigger = NewTrigger()
	return j.trigger
}
