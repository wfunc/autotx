package runner

import (
	"log"
	"os"
	"sync"

	"github.com/wfunc/autotx/task"
)

var RunnerShared *Runner

func Bootstrap() {
	if RunnerShared != nil {
		RunnerShared.All()
		log.Println("ReBootstraped")
		return
	}
	log.Println("Bootstraping...")
	RunnerShared = NewRunner()
	RunnerShared.All()
	log.Println("Bootstraped")
}

type Runner struct {
	tasks map[string]task.Task
	// ctx    context.Context
	// cancel context.CancelFunc
	wg   sync.WaitGroup
	lock sync.RWMutex
}

func NewRunner() *Runner {
	r := &Runner{}
	// ctx, cancel := context.WithCancel(context.Background())
	// r.ctx = ctx
	// r.cancel = cancel
	r.tasks = make(map[string]task.Task)
	r.wg = sync.WaitGroup{}
	r.lock = sync.RWMutex{}
	return r
}

func (r *Runner) All() {
	signTask := task.NewSignInTask("37161619", "Aa112211")
	signTask.Verbose = os.Getenv("Verbose") == "1"
	r.AddTask(signTask)
}

func (r *Runner) Loop() {
	// TODO
}

func (r *Runner) AddTask(task task.Task) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.tasks[task.TaskName()] != nil {
		log.Printf("Runner already added task with task %v", task.TaskName())
		return
	}
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		task.Run()
	}()
	r.tasks[task.TaskName()] = task
	log.Printf("Runner added task success with tasks %v", len(r.tasks))
}

func (r *Runner) Stop() int {
	r.lock.Lock()
	defer r.lock.Unlock()
	all := len(r.tasks)
	for _, task := range r.tasks {
		task.Stop()
	}
	r.wg.Wait()
	r.tasks = make(map[string]task.Task)
	return all
}
