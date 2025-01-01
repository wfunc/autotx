package task

import (
	"sync"
	"testing"
)

func TestSignIn(t *testing.T) {
	CodeURL = "https://hk51.cloudstoreapp.online/ocr?url="

	sn := NewSignInTask("60541821", "238562")
	sn.Headless = false
	sn.Verbose = true
	sn.Proxy = "http://127.0.0.1:2108"

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
