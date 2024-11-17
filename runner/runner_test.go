package runner

import (
	"testing"

	"github.com/wfunc/autotx/task"
)

func TestRunner(t *testing.T) {
	r := NewRunner()
	task1 := task.NewExampleTask("Task1")
	r.AddTask(task1)
	r.Stop()
}
