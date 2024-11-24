package task

import (
	"fmt"
	"regexp"
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

func TestRe(t *testing.T) {
	url := "steal.do?landId=3439587&from=myLand&pn=2&tuid=50765244"

	// 正则表达式匹配 tuid 的值
	re := regexp.MustCompile(`tuid=([0-9]+)`)
	match := re.FindStringSubmatch(url)

	if len(match) > 1 {
		fmt.Println("tuid:", match[1])
	} else {
		fmt.Println("tuid 不存在")
	}
}
