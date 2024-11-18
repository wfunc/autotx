package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/wfunc/autotx/conf"
	"github.com/wfunc/go/xlog"
)

type Target string

const (
	TargetQuerySeeds Target = "querySeeds"
	TargetSowSeeds   Target = "sowSeeds"
)

type FarmTask struct {
	*BaseTask
	target      Target
	Name        string
	DoneAfter   time.Duration
	successTime time.Time
}

func NewFarmTask(target Target, username, password string) *FarmTask {
	base := NewBaseTaskWithUserInfo(username, password)
	switch target {
	case TargetSowSeeds:
		base.TickerDelay = 60 * time.Second
	}
	t := &FarmTask{Name: "farm-" + string(target), BaseTask: base, target: target}
	return t
}

func (t *FarmTask) QueryShop() (shopM map[string]string, err error) {
	t.CreateChromedpContext(t.Timeout)
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
						if textContent != ">>下页" {
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

						if textContent != ">>下页" {
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

func (t *FarmTask) Run() {
	xlog.Infof("FarmTask(%v) started", t.Username)
	t.CreateChromedpContext(t.Timeout)
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
	t.Cancel()
	xlog.Infof("FarmTask(%v) done", t.Username)
}

func (t *FarmTask) Stop() {
	t.BaseTask.stop()
}

func (t *FarmTask) TaskName() string {
	return t.Username + "->" + t.Name
}

func (t *FarmTask) farm() (err error) {
	// now := time.Now()
	// if t.successTime.Year() == now.Year() && t.successTime.Month() == now.Month() && t.successTime.Day() == now.Day() {
	// 	if t.Verbose {
	// 		xlog.Infof("FarmTask(%v) sign skipped", t.Username)
	// 	}
	// 	return
	// }
	// login
	err = t.login()
	if err != nil {
		xlog.Infof("FarmTask(%v) login failed with err %v", t.Username, err)
		return
	}

	switch t.target {
	case TargetSowSeeds:
		err = t.sowSeeds()
	}

	if err == nil {
		t.successTime = time.Now()
	}
	return
}

func (t *FarmTask) sowSeeds() (err error) {
	err = chromedp.Run(t.ctx,
		chromedp.Sleep(1*time.Second),
		chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/harvestAll.do`),
		chromedp.Sleep(1*time.Second),
		chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/digLands.do`),
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			userConf := conf.Conf.GetUser(t.Username)
			setSeeds := userConf.Map("set_seeds")
			for k := range setSeeds {
				for i := 0; i < 150; i++ {
					err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/sowSeedsAll.do?seedsId=` + k).Do(ctx)
					if err != nil {
						return err
					}

					var outHTML string
					err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outHTML).Do(ctx)
					if err != nil {
						return err
					}
					seedName := setSeeds.Str(k)
					switch true {
					case strings.Contains(outHTML, "种植成功"):
						xlog.Infof("sowSeeds(%v) success with %v", t.Username, seedName)
					case strings.Contains(outHTML, "你还未开通一键"):
						t.sowSeedsManual(ctx, k)
					case strings.Contains(outHTML, "仓库中该种子已用完"):
						if strings.Contains(seedName, "金币") {
							err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/buySeeds.do?seedsId=` + k + "&num=1").Do(ctx)
							if err != nil {
								return err
							}
							err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outHTML).Do(ctx)
							if err != nil {
								return err
							}
							if strings.Contains(outHTML, "恭喜,你成功购买了") {
								err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/sowSeedsAll.do?seedsId=` + k).Do(ctx)
								if err != nil {
									return err
								}
								err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outHTML).Do(ctx)
								if err != nil {
									return err
								}
								if strings.Contains(outHTML, "种植成功") || strings.Contains(outHTML, "成功种植") {
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
						var outHTMLErr string
						err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc > div.mainarea > div.dotline0`, &outHTMLErr).Do(ctx)
						if err != nil {
							return err
						}
						if !strings.Contains(outHTMLErr, "还没有可种植的土地") {
							// xlog.Infof("[%v]种菜失败☹️[%v]原因：%v", account, v, outHTMLErr)
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
	return
}

func (t *FarmTask) sowSeedsManual(ctx context.Context, seedID string) (err error) {
	err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/index.do`).Do(ctx)
	if err != nil {
		return
	}

	var needDoM = map[string]string{}
	var outHtml string
	err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outHtml).Do(ctx)
	if err != nil {
		return
	}

	err = ExtractLinksWithPrefix(outHtml, needDoM, "https://tx.com.cn/plugins/farm/cs/", []string{"myBag.do", "digLand.do"})
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
			var outHtml string
			err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outHtml).Do(ctx)
			if err != nil {
				return
			}
			if strings.Contains(outHtml, "种植成功") {
				xlog.Infof("sowSeedsManual(%v) success with %v", t.Username, seedID)
			} else {
				xlog.Infof("sowSeedsManual(%v) failed with %v html %v", t.Username, seedID, outHtml)
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
	ExtractLinksWithPrefix(outerHTML, m, "https://tx.com.cn/plugins/farm/cs/", []string{"seedsInfo.do"})
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

	var outHTML string
	err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outHTML).Do(ctx)
	if err != nil {
		return err
	}
	if strings.Contains(outHTML, "种植成功") {
		xlog.Infof("sowSeedsRetry(%v) success with %v", t.Username, seedID)
	} else {
		var outHTMLErr string
		err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc > div.mainarea > div.dotline0`, &outHTMLErr).Do(ctx)
		if err != nil {
			return err
		}
		if !strings.Contains(outHTMLErr, "还没有可种植的土地") {
		}
	}

	if strings.Contains(outHTML, "仓库中该种子已用完") {
		if strings.Contains(seedName, "金币") {
			// 购买种子重新种一次
			err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/buySeeds.do?seedsId=` + seedID + "&num=1").Do(ctx)
			if err != nil {
				// xlog.Infof("[%v]种菜[%v]错误:%v", account, v, err)
				return err
			}
			// 恭喜,你成功购买了
			err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outHTML).Do(ctx)
			if err != nil {
				return err
			}
			if strings.Contains(outHTML, "恭喜,你成功购买了") {
				err = chromedp.Navigate(`https://tx.com.cn/plugins/farm/cs/sowSeedsAll.do?seedsId=` + seedID).Do(ctx)
				if err != nil {
					return err
				}

				err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outHTML).Do(ctx)
				if err != nil {
					return err
				}
				if strings.Contains(outHTML, "种植成功") || strings.Contains(outHTML, "成功种植") {
					xlog.Infof("sowSeedsRetry(%v) success with %v", t.Username, seedName)
				} else {
					return nil
				}
			}
		}

	}
	return
}
