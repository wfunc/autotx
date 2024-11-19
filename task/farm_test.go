package task

import (
	"fmt"
	"testing"
)

func TestFarmGetTime(t *testing.T) {
	fram := NewFarmTask(TargetSowSeeds, "37161619", "Aa112211")
	fram.Verbose = true
	fram.Headless = false
	fram.CreateChromedpContext(fram.Timeout)

	fmt.Println(fram.GetTime())
	fram.Cancel()
}
