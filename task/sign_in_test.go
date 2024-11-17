package task

import (
	"testing"
	"time"
)

func TestSignIn(t *testing.T) {
	sn := NewSignInTask("37161619", "Aa112211")
	sn.Headless = false
	sn.Verbose = true
	go sn.Run()
	time.Sleep(15 * time.Second)
	sn.Stop()
	time.Sleep(3 * time.Second)
	go sn.Run()
	sn.Stop()
}
