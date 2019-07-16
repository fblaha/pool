# pool 
[![GoDoc](https://godoc.org/github.com/fblaha/pool?status.svg)](https://godoc.org/github.com/fblaha/pool)
[![Build Status](https://api.travis-ci.org/fblaha/pool.svg?branch=master)](https://api.travis-ci.org/fblaha/pool)
[![Sourcegraph](https://sourcegraph.com/github.com/fblaha/pool/-/badge.svg)](https://sourcegraph.com/github.com/fblaha/pool?badge)

This library/module provides goroutine pool executor with the completion wait support. 
An executor that executes each submitted job using one of possibly several pooled goroutines.
An user has 2 options how to implement jobs:

1. By type implementing Worker interface (`executor.Submit`) 
```go
type Worker interface {
	Work()
}
```

2. By pure function (`executor.SubmitFunc`)
```go
type WorkerFunc func()
```


Download:
```shell
go get github.com/fblaha/pool
```

* * *

Simple Example:
```go
// creates a new executor with a pool of 10 goroutines
executor := NewExecutor(10)

// terminates worker go routines and frees allocated resources
defer executor.ShutdownGracefully()

// submits 2 worker functions simulating time consuming jobs
executor.SubmitFunc(
    func() {
        time.Sleep(100 * time.Millisecond)
        fmt.Println("hard work done")
    },
    func() {
        time.Sleep(200 * time.Millisecond)
        fmt.Println("double hard work done")

    })

// waits for the completion of above submitted work
// the wait call influences the order of the output lines
executor.Wait()

// submits a one more worker function simulating a short time job
executor.SubmitFunc(
    func() {
        // this line should occur in the output last due to above wait call
        fmt.Println("easy work done")
    })

// Output:
//hard work done
//double hard work done
//easy work done
```

Real Example:
 * godoc example https://godoc.org/github.com/fblaha/pool
 * demonstrates output channel for collecting results
 * demonstrates errors propagation via output channel
