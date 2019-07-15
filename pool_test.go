package pool

import (
	"sync"
	"testing"
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
