package task

import (
	"testing"
)

func TestExampleTask(t *testing.T) {
	task := NewExampleTask("ExampleTask")
	// task.Headless = false
	// task.DoneAfter = 3 * time.Second
	task.Run()
	task.stop()
	task.Run()
}
