package runner

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/wfunc/autotx/conf"
	"github.com/wfunc/autotx/task"
	"github.com/wfunc/go/xlog"
	"github.com/wfunc/util/xmap"
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
	users := conf.Conf.GetUsers()
	if len(users) < 1 {
		return
	}
	seed := len(conf.Conf.GetSeedsRevert()["1"]) < 1
	for username, userConf := range users {
		password := userConf.Str("password")
		if seed {
			farmTask := task.NewFarmTask(task.TargetQuerySeeds, username, password)
			farmTask.Verbose = os.Getenv("Verbose") == "1"
			farmTask.CreateChromedpContext(farmTask.Timeout)
			seedM, _ := farmTask.QueryShop()
			conf.Conf.SetSeeds(seedM)
			if len(seedM) > 0 {
				seed = false
			}
			seedM, _ = farmTask.PaySeeds()
			conf.Conf.SetSeeds(seedM)
			farmTask.Cancel()
		}
		r.StartTask(username, password)
	}
}

func (r *Runner) Reload(username string) {
	userConf := conf.Conf.GetUser(username)
	password := userConf.Str("password")
	r.StartTask(username, password)
}

func (r *Runner) StopTask(username string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, task := range r.tasks {
		if strings.Contains(task.TaskName(), username) {
			task.Stop()
			delete(r.tasks, task.TaskName())
		}
	}
}

func (r *Runner) Loop() {
	// TODO
}

func (r *Runner) ListTask() (result xmap.M) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	result = xmap.M{}
	for _, task := range r.tasks {
		result[task.TaskName()] = task.Info()
	}
	return
}

func (r *Runner) AddTask(task task.Task) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.tasks[task.TaskName()] != nil {
		xlog.Infof("Runner already added task with task %v", task.TaskName())
		return
	}
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		task.Run()
	}()
	r.tasks[task.TaskName()] = task
	xlog.Infof("Runner added task success with tasks %v", len(r.tasks))
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

func (r *Runner) StartTask(username, password string) {
	// sign
	signTask := task.NewSignInTask(username, password)
	signTask.Verbose = os.Getenv("Verbose") == "1"
	r.AddTask(signTask)

	// exchange card
	exchangeCardTask := task.NewExchangeCardTask(username, password)
	exchangeCardTask.Verbose = os.Getenv("Verbose") == "1"
	r.AddTask(exchangeCardTask)

	// sow seeds
	sowSeeds := task.NewFarmTask(task.TargetSowSeeds, username, password)
	sowSeeds.Verbose = os.Getenv("Verbose") == "1"
	r.AddTask(sowSeeds)

	// farm water
	water := task.NewFarmTask(task.TargetWater, username, password)
	water.Verbose = os.Getenv("Verbose") == "1"
	r.AddTask(water)

}
