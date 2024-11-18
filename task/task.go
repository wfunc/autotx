package task

import "github.com/wfunc/util/xmap"

type Task interface {
	TaskName() string
	Run()
	Stop()
	Info() xmap.M
}
