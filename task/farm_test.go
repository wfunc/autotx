package task

import (
	"fmt"
	"testing"
	"time"
)

func TestFarmGetTime(t *testing.T) {
	farm := NewFarmTask(TargetSowSeeds, "37161619", "Aa112211")
	farm.Verbose = true
	farm.Headless = false
	farm.CreateChromedpContext(farm.Timeout)

	fmt.Println(farm.GetTime())
	farm.Cancel()
}

func TestFarmWater(t *testing.T) {
	farm := NewFarmTask(TargetWater, "37161619", "Aa112211")
	// farm.Verbose = true
	farm.Headless = false
	farm.Timeout = 60 * time.Second
	farm.CreateChromedpContext(farm.Timeout)
	farm.login()
	farm.water()
	fmt.Println("water done")
	farm.Cancel()
}
