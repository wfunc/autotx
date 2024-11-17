package task

import "os"

var (
	CodeAPI = true
	CodeURL = ""
)

func BootstrapConfig() {
	CodeURL = os.Getenv("CodeURL")
}

const (
	canceled = "context canceled"
	deadline = "context deadline exceeded"
)
