package pool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestExecutorStartShutdown(t *testing.T) {
	executor := NewExecutor(10)
	executor.ShutdownGracefully()
}

type mockWork struct {
	sync.WaitGroup
}

func (w *mockWork) Work() {
	w.Done()
}

func TestExecutorSubmit(t *testing.T) {
	executor := NewExecutor(10)
	defer executor.ShutdownGracefully()
	var work mockWork
	work.Add(100)
	for i := 0; i < 100; i++ {
		executor.Submit(&work)
	}
	work.Wait()
}

func ExampleExecutor_SubmitFunc() {
	// creates a new executor with a pool of 10 goroutines
	executor := NewExecutor(10)

	// terminates worker go routines and frees resources
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
}

func TestExecutorSubmitFunc(t *testing.T) {
	executor := NewExecutor(10)
	var wg sync.WaitGroup
	wg.Add(3)
	fc := func() {
		wg.Done()
	}
	executor.SubmitFunc(fc, fc, fc)
	executor.Wait()
	wg.Wait()
	defer executor.ShutdownGracefully()
}
