package task

import (
	"testing"
)

func TestLogin(t *testing.T) {
	base := NewBaseTaskWithUserInfo("37161619", "Aa112211")
	base.Headless = false
	base.CreateChromedpContext(base.Timeout)
	var err = base.login()
	if err != nil {
		t.Error(err)
	}
	base.stop()
	err = base.login()
	if err == nil {
		t.Error("login success")
	}
}
