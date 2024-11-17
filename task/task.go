package task

type Task interface {
	TaskName() string
	Run()
	Stop()
}
