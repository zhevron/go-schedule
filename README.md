go-schedule - Task scheduler library
====================================

[![wercker status](https://app.wercker.com/status/d9189749c854d947e76b9a4a675ac7f9/s "wercker status")](https://app.wercker.com/project/bykey/d9189749c854d947e76b9a4a675ac7f9)
[![Coverage Status](https://coveralls.io/repos/zhevron/go-schedule/badge.svg?branch=HEAD)](https://coveralls.io/r/zhevron/go-schedule)
[![GoDoc](https://godoc.org/gopkg.in/zhevron/go-schedule.v0/schedule?status.svg)](https://godoc.org/gopkg.in/zhevron/go-schedule.v0/schedule)

**go-schedule** is a task scheduling library for [Go](https://golang.org/).  

For package documentation, refer to the GoDoc badge above.

## Installation

```
go get gopkg.in/zhevron/go-schedule.v0/schedule
```

## Usage

```go
package main

import (
  "fmt"

  "gopkg.in/zhevron/go-schedule.v0/schedule"
)

func MyJob(scheduler *schedule.Scheduler) {
  fmt.Println("Running MyJob")
  scheduler.Stop()
}

func main() {
  scheduler := schedule.NewScheduler()

  job := schedule.NewJob("MyJob", MyJob, scheduler);
  job.Schedule().Every("15m").From(time.Now())

  scheduler.Add(job)
  scheduler.Start()

  for scheduler.Running() {
    switch {
    case res := <-scheduler.Results():
      fmt.Printf("job %#s returned result: %v\n", res.Name, res.Results)

    case err := <-scheduler.Errors():
      fmt.Printf("job %#s encountered an error: %s\n", err.Name, err.Error)
    }
  }
}
```

## License

**go-schedule** is licensed under the [MIT license](http://opensource.org/licenses/MIT).
