package task

import (
	"sync"
	"testing"
)

func TestSignIn(t *testing.T) {
	CodeURL = "https://ocr.rosetts.com/ocr?url="

	sn := NewSignInTask("37161619", "Aa112211")
	sn.Headless = false
	sn.Verbose = true

	wg := sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sn.Run()
		}()
		sn.Stop()
		wg.Wait()
	}
}
