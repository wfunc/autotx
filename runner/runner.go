package runner

import (
	"context"
	"log"
	"sync"

	"github.com/wfunc/autotx/task"
)

type Runner struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	lock   sync.RWMutex
}

func NewRunner() *Runner {
	r := &Runner{}
	ctx, cancel := context.WithCancel(context.Background())
	r.ctx = ctx
	r.cancel = cancel
	r.wg = sync.WaitGroup{}
	r.lock = sync.RWMutex{}
	return r
}

func (r *Runner) Run() {
	r.wg.Wait()
}

func (r *Runner) AddTask(task task.Task) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		if err := task.Run(r.ctx); err != nil {
			log.Printf("Task error: %v", err)
		}
	}()
}

func (r *Runner) Stop() {
	r.cancel()
	r.wg.Wait()
}
