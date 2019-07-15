package pool

import (
	"sync"
)

// Worker does the work
type Worker interface {
	Work()
}

// WorkerFunc does the work
type WorkerFunc func()

// PoolExecutor distributes works to the pool of goroutines
type PoolExecutor struct {
	// tracks completion of submitted work
	workWG sync.WaitGroup
	// tracks pool goroutines which process the incoming work
	poolWG sync.WaitGroup
	// incoming work
	todo chan WorkerFunc
}

// NewPoolExecutor constructor
func NewPoolExecutor(concurrency int) *PoolExecutor {
	todo := make(chan WorkerFunc)
	executor := PoolExecutor{todo: todo}
	for i := 0; i < concurrency; i++ {
		executor.poolWG.Add(1)
		go executor.handleWork()
	}
	return &executor
}

func (e *PoolExecutor) handleWork() {
	defer e.poolWG.Done()
	for w := range e.todo {
		w()
		e.workWG.Done()
	}
}

// SubmitFunc submits the work for execution
func (e *PoolExecutor) SubmitFunc(workers ...WorkerFunc) {
	for _, w := range workers {
		// ensures that shutdown waits for completion of submitted work
		e.workWG.Add(1)
		// submits work
		e.todo <- w
	}
}

// Submit submits the work for execution
func (e *PoolExecutor) Submit(workers ...Worker) {
	for _, w := range workers {
		e.SubmitFunc(w.Work)
	}
}

// ShutdownGracefully waits for completion of the submitted work and terminates worker goroutines
// and frees allocated resources. The executor can no longer be used after this call.
func (e *PoolExecutor) ShutdownGracefully() {
	// waits for completion of submitted work
	e.Wait()
	close(e.todo)
	// waits for completion of the pool goroutines
	e.poolWG.Wait()
}

// Wait waits for completion of the submitted work
func (e *PoolExecutor) Wait() {
	e.workWG.Wait()
}
