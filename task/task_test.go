package task

import (
	"context"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	task := &ExampleTask{Name: "Task1"}
	ctx := context.Background()
	task.Run(ctx)
	time.Sleep(3 * time.Second)
	task.Stop()
}
