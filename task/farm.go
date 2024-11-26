package task

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/wfunc/autotx/conf"
	"github.com/wfunc/go/xlog"
	"github.com/wfunc/util/xmap"
)

type Target string

const (
	TargetQuerySeeds Target = "querySeeds"
	TargetSowSeeds   Target = "sowSeeds"
	TargetWater      Target = "water"
)

type FarmTask struct {
	*BaseTask
	target      Target
	Name        string
	DoneAfter   time.Duration
	nextTime    time.Time
	successTime time.Time
}

func NewFarmTask(target Target, username, password string) *FarmTask {
	base := NewBaseTaskWithUserInfo(username, password)
	switch target {
	case TargetSowSeeds:
		base.TickerDelay = 60 * time.Second
	case TargetWater:
		base.TickerDelay = 30 * time.Second
	}
	base.Timeout = 10 * time.Minute
	t := &FarmTask{Name: "farm-" + string(target), BaseTask: base, target: target}
	base.ParentTask = t
	return t
}

func (t *FarmTask) QueryShop() (shopM map[string]string, err error) {
	t.login()
	shopM = map[string]string{}
	err = chromedp.Run(t.ctx,
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < 2; i++ {
				var err = chromedp.Navigate(fmt.Sprintf(`https://tx.com.cn/plugins/farm/cs/shop.do?category=1&lv=%v`, i)).Do(ctx)
				if err != nil {
					if t.Verbose {
						xlog.Infof("FarmTask(%v) query shop failed with err %v", t.Username, err)
					}
					return err
				}
				var outerHTML string
				err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
				if err == nil {
					ParseShopHTML(outerHTML, shopM)
					var textContent, href string
					for {
						err = chromedp.Evaluate(`document.querySelector('body > div.mainareaOutside_pc > div.mainareaCenter_pc > div:nth-child(13) > a').textContent`, &textContent).Do(ctx)
						if err != nil {
							break
						}
						err = chromedp.Evaluate(`document.querySelector('body > div.mainareaOutside_pc > div.mainareaCenter_pc > div:nth-child(13) > a').href`, &href).Do(ctx)
						if err != nil {
							break
						}
						if textContent != "下页" {
							break
						}
						chromedp.Navigate(href).Do(ctx)
						err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
						if err == nil {
							ParseShopHTML(outerHTML, shopM)
						}
					}
				}

			}
			for i := 0; i < 4; i++ {
				chromedp.Navigate(fmt.Sprintf(`https://tx.com.cn/plugins/farm/cs/shop.do?category=0&lv=%v`, i)).Do(ctx)
				var outerHTML string
				err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
				if err == nil {
					ParseShopHTML(outerHTML, shopM)
					var textContent, href string
					for {
						err = chromedp.Evaluate(`document.querySelector('body > div.mainareaOutside_pc > div.mainareaCenter_pc > div:nth-child(13) > a').textContent`, &textContent).Do(ctx)
						if err != nil {
							break
						}
						err = chromedp.Evaluate(`document.querySelector('body > div.mainareaOutside_pc > div.mainareaCenter_pc > div:nth-child(13) > a').href`, &href).Do(ctx)
						if err != nil {
							break
						}
						// xlog.Infof("textContent is %v, href is %v", textContent, href)

						if textContent != "下页" {
							break
						}
						chromedp.Navigate(href).Do(ctx)
						err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
						if err == nil {
							ParseShopHTML(outerHTML, shopM)
						}
					}
				}

			}
			return nil
		}),
	)
	return
}

func (t *FarmTask) PaySeeds() (shopM map[string]string, err error) {
	shopM = map[string]string{}
	err = chromedp.Run(t.ctx,
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// for i := 0; i < 4; i++ {
			chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/paySeeds.do`).Do(ctx)
			var outerHTML string
			err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
			if err == nil {
				ParseShopHTML(outerHTML, shopM)
				var textContent, href string
				for {
					err = chromedp.Evaluate(`document.querySelector('body > div.mainareaOutside_pc > div.mainareaCenter_pc > div:nth-child(14) > a').textContent`, &textContent).Do(ctx)
					if err != nil {
						break
					}
					err = chromedp.Evaluate(`document.querySelector('body > div.mainareaOutside_pc > div.mainareaCenter_pc > div:nth-child(14) > a').href`, &href).Do(ctx)
					if err != nil {
						break
					}
					// xlog.Infof("textContent is %v, href is %v", textContent, href)

					if textContent != "下页" {
						break
					}
					chromedp.Navigate(href).Do(ctx)
					err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
					if err == nil {
						ParseShopHTML(outerHTML, shopM)
					}
				}
			}

			// }
			return nil
		}),
	)
	// xlog.Infof("shops is %v", converter.JSON(shopM))
	return
}

func (t *FarmTask) Run() {
	xlog.Infof("FarmTask(%v) started", t.Username)
	t.nextTime = time.Time{}
	t.farm()
	ticker := time.NewTicker(t.TickerDelay)
	defer ticker.Stop()
	running := true
	for running {
		select {
		case <-t.exiter:
			running = false
		case <-ticker.C:
			t.farm()
		}
	}
	if t.Cancel != nil {
		t.Cancel()
	}
	xlog.Infof("FarmTask(%v) done", t.Username)
}

func (t *FarmTask) Stop() {
	t.BaseTask.stop()
}

func (t *FarmTask) Info() (result xmap.M) {
	result = t.BaseInfo()
	result["success_time"] = t.successTime
	return
}

func (t *FarmTask) TaskName() string {
	return t.Username + "->" + t.Name
}

func (t *FarmTask) farm() (err error) {
	// check next time
	now := time.Now()
	if t.nextTime.Sub(now) > 0 {
		if t.Verbose {
			xlog.Infof("FarmTask(%v) next time is %v, now is %v", t.Username, t.nextTime, now)
		}
		return
	}
	userConf := conf.Conf.GetUser(t.Username)
	setSeeds := userConf.Map("set_seeds")
	// check user set seeds
	if setSeeds.Length() < 1 {
		if t.Verbose {
			xlog.Infof("FarmTask(%v) set_seeds is empty", t.Username)
		}
		return
	}
	t.CreateChromedpContext(t.Timeout)
	defer t.Cancel()
	err = t.login()
	if err != nil {
		xlog.Infof("FarmTask(%v) login failed with err %v", t.Username, err)
		return
	}

	switch t.target {
	case TargetSowSeeds:
		err = t.sowSeeds(setSeeds)
		t.nextTime = time.Now().Add(t.GetTime())
		xlog.Infof("FarmTask(%v) will do nextime on %v", t.TaskName(), t.nextTime.Format("2006-01-02 15:04:05"))
	case TargetWater:
		err = t.water()
	}

	if err == nil {
		t.successTime = time.Now()
	}
	return
}

func (t *FarmTask) sowSeeds(setSeeds xmap.M) (err error) {
	err = chromedp.Run(t.ctx,
		chromedp.Sleep(1*time.Second),
		chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/harvestAll.do`),
		chromedp.Sleep(1*time.Second),
		chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/digLands.do`),
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {

			for k := range setSeeds {
				for i := 0; i < 150; i++ {
					err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/sowSeedsAll.do?seedsId=` + k).Do(ctx)
					if err != nil {
						return err
					}
					var outerHTML string
					err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
					if err != nil {
						return err
					}
					seedName := setSeeds.Str(k)
					switch true {
					case strings.Contains(outerHTML, "种植成功"):
						xlog.Infof("sowSeeds(%v) success with %v", t.Username, seedName)
					case strings.Contains(outerHTML, "你还未开通一键"):
						t.sowSeedsManual(ctx, k)
					case strings.Contains(outerHTML, "仓库中该种子已用完"):
						if strings.Contains(seedName, "金币") {
							err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/buySeeds.do?seedsId=` + k + "&num=1").Do(ctx)
							if err != nil {
								return err
							}
							err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
							if err != nil {
								return err
							}
							if strings.Contains(outerHTML, "恭喜,你成功购买了") {
								err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/sowSeedsAll.do?seedsId=` + k).Do(ctx)
								if err != nil {
									return err
								}
								err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
								if err != nil {
									return err
								}
								if strings.Contains(outerHTML, "种植成功") || strings.Contains(outerHTML, "成功种植") {
									xlog.Infof("sowSeeds(%v) success with %v", t.Username, seedName)
								} else {
									return nil
								}
							} else {
								err = t.sowSeedsLevel(ctx)
								if err != nil {
									return err
								}
							}
						}
					default:
						var outerHTMLErr string
						err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc > div.mainarea > div.dotline0`, &outerHTMLErr).Do(ctx)
						if err != nil {
							return err
						}
						if !strings.Contains(outerHTMLErr, "还没有可种植的土地") {
							// xlog.Infof("[%v]种菜失败☹️[%v]原因：%v", account, v, outerHTMLErr)
						}

						err = t.sowSeedsLevel(ctx)
						if err != nil {
							return err
						}
					}
				}
			}

			return nil
		}),
	)
	if err == nil {
		conf.Conf.UpdateUser(t.Username, "sowSeeds", time.Now().Format(`2006-01-02 15:04:05`))
	}
	return
}

func (t *FarmTask) sowSeedsManual(ctx context.Context, seedID string) (err error) {
	err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/index.do`).Do(ctx)
	if err != nil {
		return
	}

	var needDoM = map[string]string{}
	var outerHTML string
	err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
	if err != nil {
		return
	}

	_, err = ExtractLinksWithPrefix(outerHTML, &needDoM, "https://tx.com.cn/plugins/farm/cs/", []string{"myBag.do", "digLand.do"})
	if err != nil {
		xlog.Infof("sowSeedsManual(%v) failed with err %v", t.Username, err)
		return
	}
	for k, v := range needDoM {
		if strings.Contains(k, "z=") {
			continue
		}
		err = chromedp.Navigate(v + k).Do(ctx)
		if err != nil {
			return
		}
		if strings.Contains(k, "myBag.do") {
			landID := strings.Split(k, "myBag.do?")[1]
			err = chromedp.Navigate(fmt.Sprintf(`https://tx.com.cn/plugins/farm/sowSeeds.do?landId=%v&seedsId=%v`, landID, seedID)).Do(ctx)
			if err != nil {
				return
			}
			var outerHTML string
			err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
			if err != nil {
				return
			}
			if strings.Contains(outerHTML, "种植成功") {
				xlog.Infof("sowSeedsManual(%v) success with %v", t.Username, seedID)
			} else {
				xlog.Infof("sowSeedsManual(%v) failed with %v html %v", t.Username, seedID, outerHTML)
			}
		}
	}
	return
}

func (t *FarmTask) sowSeedsLevel(ctx context.Context) (err error) {
	err = chromedp.Navigate("https://tx.com.cn/plugins/farm/cs/myInfo.do").Do(ctx)
	if err != nil {
		return err
	}
	var outerHTML string
	err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
	if err != nil {
		return err
	}
	m := map[string]string{}
	ExtractLinksWithPrefix(outerHTML, &m, "https://tx.com.cn/plugins/farm/cs/", []string{"seedsInfo.do"})
	seedM := conf.Conf.GetSeedsRevert()
	for k := range m {
		seedsID := strings.Split(k, "seedsId=")[1]
		err = t.sowSeedsRetry(ctx, seedsID, seedM[seedsID])
		if err != nil {
			return err
		}
	}
	return
}

func (t *FarmTask) sowSeedsRetry(ctx context.Context, seedID, seedName string) (err error) {
	err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/sowSeedsAll.do?seedsId=` + seedID).Do(ctx)
	if err != nil {
		return err
	}

	var outerHTML string
	err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
	if err != nil {
		return err
	}
	if strings.Contains(outerHTML, "种植成功") {
		xlog.Infof("sowSeedsRetry(%v) success with %v", t.Username, seedID)
	} else {
		var outerHTMLErr string
		err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc > div.mainarea > div.dotline0`, &outerHTMLErr).Do(ctx)
		if err != nil {
			return err
		}
		if !strings.Contains(outerHTMLErr, "还没有可种植的土地") {
		}
	}

	if strings.Contains(outerHTML, "仓库中该种子已用完") {
		if strings.Contains(seedName, "金币") {
			// 购买种子重新种一次
			err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/buySeeds.do?seedsId=` + seedID + "&num=1").Do(ctx)
			if err != nil {
				// xlog.Infof("[%v]种菜[%v]错误:%v", account, v, err)
				return err
			}
			// 恭喜,你成功购买了
			err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
			if err != nil {
				return err
			}
			if strings.Contains(outerHTML, "恭喜,你成功购买了") {
				err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/sowSeedsAll.do?seedsId=` + seedID).Do(ctx)
				if err != nil {
					return err
				}

				err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
				if err != nil {
					return err
				}
				if strings.Contains(outerHTML, "种植成功") || strings.Contains(outerHTML, "成功种植") {
					xlog.Infof("sowSeedsRetry(%v) success with %v", t.Username, seedName)
				} else {
					return nil
				}
			}
		}

	}
	return
}

func (t *FarmTask) GetTime() (minTime time.Duration) {
	t.login()
	minTime = time.Duration(10 * time.Hour)
	err := chromedp.Run(t.ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			url := "https://tx.com.cn/plugins/farm/index.do?pn=1"
			var err = chromedp.Navigate(url).Do(ctx)
			if err != nil {
				return err
			}
			var outerHTML string
			err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
			if err == nil {
				minTime = t.extractTimes(outerHTML, minTime)
				for {
					selector := `//a[contains(text(), "下页")]`
					if !t.clickNext(ctx, selector) {
						break
					}
					err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
					if err == nil {
						t.extractTimes(outerHTML, minTime)
					}
				}
			}

			// }
			return nil
		}),
	)
	if err != nil {
		return
	}

	return
}

// func extractTimes(html string) {
// 	// 定义正则表达式匹配时间信息
// 	re := regexp.MustCompile(`(\d+)小时(\d+)分钟`)
// 	matches := re.FindAllStringSubmatch(html, -1)

// 	// 输出提取的时间
// 	for _, match := range matches {
// 		if len(match) == 3 {
// 			fmt.Printf("提取的时间: %s小时%s分钟\n", match[1], match[2])
// 		}
// 	}
// 	time.Duration
// }

func (t *FarmTask) extractTimes(html string, minTime time.Duration) (currentMin time.Duration) {
	re := regexp.MustCompile(`(\d+)小时(\d+)分钟`)
	matches := re.FindAllStringSubmatch(html, -1)
	currentMin = minTime
	for _, match := range matches {
		if len(match) == 3 {
			hours, minutes := match[1], match[2]
			hourDuration, err1 := time.ParseDuration(fmt.Sprintf("%sh", hours))
			minuteDuration, err2 := time.ParseDuration(fmt.Sprintf("%sm", minutes))
			if err1 != nil || err2 != nil {
				continue
			}
			totalDuration := hourDuration + minuteDuration
			if t.Verbose {
				xlog.Infof("extractTimes ---> %v", totalDuration)
			}
			if totalDuration < currentMin {
				currentMin = totalDuration
			}
		}
	}
	return
}
