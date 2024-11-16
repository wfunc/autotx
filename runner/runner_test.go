package runner

import (
	"testing"
	"time"

	"github.com/wfunc/autotx/task"
)

func TestRunner(t *testing.T) {
	r := NewRunner()
	r.Run()
	// 添加任务
	task1 := &task.ExampleTask{Name: "Task1"}
	r.AddTask(task1)
	// r.Stop()
	time.Sleep(3 * time.Second)
	r.Stop()
}
