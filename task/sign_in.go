package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/wfunc/autotx/conf"
	"github.com/wfunc/go/xlog"
	"github.com/wfunc/util/xmap"
)

type SignInTask struct {
	*BaseTask
	Name        string
	DoneAfter   time.Duration
	successTime time.Time
	failed      int
}

func NewSignInTask(username, password string) *SignInTask {
	base := NewBaseTaskWithUserInfo(username, password)
	base.Timeout = 2 * time.Minute
	t := &SignInTask{Name: "sign-in", BaseTask: base}
	base.ParentTask = t
	return t
}

func (t *SignInTask) Run() {
	t.UserAgent = IphoneUserAgent
	user := conf.Conf.GetUser(t.Username)
	signIN := user.Str("signIN")
	if len(signIN) > 0 {
		layout := "2006-01-02 15:04:05"

		// 使用 time.Parse 将字符串解析为 time.Time
		parsedTime, err := time.ParseInLocation(layout, signIN, time.Local)
		if err == nil {
			t.successTime = parsedTime
			xlog.Infof("SignInTask(%v) signIN time is %v", t.Username, parsedTime)
		}
	} else {
		t.clear()
	}

	xlog.Infof("SignInTask(%v) started", t.Username)

	t.signIN()
	ticker := time.NewTicker(t.TickerDelay)
	defer ticker.Stop()
	running := true
	for running {
		select {
		case <-t.exiter:
			running = false
		case <-ticker.C:
			result, _ := t.signIN()
			if len(result) > 0 {
				t.failed++
				xlog.Infof("SignInTask(%v) sign failed(%v) with %v will sleep on %v", t.Username, t.failed, result, t.Timeout)
				time.Sleep(t.Timeout)
			}
		}
	}
	if t.Cancel != nil {
		t.Cancel()
	}
	xlog.Infof("SignInTask(%v) done with %v", t.Username, t.Cancel)
}

func (t *SignInTask) Stop() {
	t.BaseTask.stop()
}

func (t *SignInTask) Info() (result xmap.M) {
	result = t.BaseInfo()
	result["success_time"] = t.successTime
	return
}

func (t *SignInTask) TaskName() string {
	return t.Username + "->" + t.Name
}

func (t *SignInTask) clear() {
	t.successTime = time.Time{}
	t.failed = 0
}

func (t *SignInTask) signIN() (result string, err error) {
	now := time.Now()
	makeTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 6, 0, 0, now.Location())
	if t.successTime.Year() == now.Year() && t.successTime.Month() == now.Month() && t.successTime.Day() == now.Day() {
		if t.Verbose {
			xlog.Infof("SignInTask(%v) skipped and next will after on %v", t.Username, makeTime.AddDate(0, 0, 1).Sub(now))
		}
		if makeTime.AddDate(0, 0, 1).Sub(now) > time.Hour {
			time.Sleep(300 * time.Second)
		}
		return
	}
	// 如果now小于makeTime,直到makeTime到了再继续
	if now.Before(makeTime) {
		xlog.Infof("SignInTask(%v) sign will on %v start", t.Username, makeTime)
		time.Sleep(makeTime.Sub(now))
	}
	lock := ChromeManagerInstance.GetUserLock(t.Username)
	lock.Lock()
	defer lock.Unlock()
	t.CreateChromedpContext(t.Timeout)
	defer t.Cancel()

	// if len(t.Proxy) > 0 {
	// 	err = chromedp.Run(t.ctx, chromedp.Navigate(`https://ip.cn`), chromedp.Sleep(5*time.Second))
	// 	if err != nil {
	// 		xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
	// 		return
	// 	}
	// }
	// login
	err = t.login()
	if err != nil {
		xlog.Infof("SignInTask(%v) login failed with err %v", t.Username, err)
		return
	}
	// sign
	err = chromedp.Run(t.ctx,
		chromedp.Sleep(1*time.Second),
		chromedp.Navigate(`https://tx.com.cn/space/index.do`),
		chromedp.Sleep(1*time.Second),
		chromedp.Navigate(`https://tx.com.cn/plugins/pet2/cs/exchange.do?batch=1&act=1`),
		chromedp.Sleep(1*time.Second),
		chromedp.Navigate(`https://tx.com.cn/plugins/pet2/cs/exchange.do?batch=1&act=3`),
		chromedp.Sleep(1*time.Second),
		chromedp.Navigate(`https://tx.com.cn/plugins/pet2/cs/exchange.do?batch=1&act=5`),
		chromedp.Sleep(1*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			defer func() {
				syndic := conf.Conf.LoadDo("syndic")
				if len(syndic) > 0 {
					friendIDs := strings.Split(syndic, ",")
					for _, friendID := range friendIDs {
						if t.Username == friendID {
							continue
						}
						err = chromedp.Navigate(fmt.Sprintf("https://tx.com.cn/show/syndic.do?score=1&friendId=%v", friendID)).Do(ctx)
						if err != nil {
							xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
							return
						}
						if t.Verbose {
							xlog.Infof("SignInTask(%v) syndic success with %v", t.Username, friendID)
						}
						err = chromedp.Sleep(1 * time.Second).Do(ctx)
						if err != nil {
							xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
							return
						}
					}
				}
			}()

			if t.Verbose {
				xlog.Infof("SignInTask(%v) sign start--->1", t.Username)
			}

			err = chromedp.Navigate(`https://tx.com.cn/activity/qq/cs/sign.do`).Do(ctx)
			if err != nil {
				xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
				return err
			}
			err = chromedp.Sleep(1 * time.Second).Do(ctx)
			if err != nil {
				xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
				return err
			}
			if t.Verbose {
				xlog.Infof("SignInTask(%v) sign start--->2", t.Username)
			}
			var str string
			err = chromedp.OuterHTML(`body`, &str).Do(ctx)
			if err != nil {
				xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
				return err
			}
			if t.Verbose {
				xlog.Infof("SignInTask(%v) sign start--->3", t.Username)
			}
			if strings.Contains(str, "今日已签到过") {
				return nil
			}
			switch true {
			case strings.Contains(str, "请输入图片中的验证码"):
				if t.Verbose {
					xlog.Infof("SignInTask(%v) sign start--->4", t.Username)
				}
				err = chromedp.Evaluate(`document.querySelector("body > form > img").src`, &str).Do(ctx)
				if err != nil {
					time.Sleep(300 * time.Second)
					xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}

				if t.Verbose {
					xlog.Infof("SignInTask(%v) sign start--->5", t.Username)
				}
				authnum := getCode(str)
				err = chromedp.Sleep(1 * time.Second).Do(ctx)
				if err != nil {
					xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}
				if t.Verbose {
					xlog.Infof("SignInTask(%v) sign start---> url %v code %v", t.Username, str, authnum)
				}
				err = chromedp.WaitVisible(`body > form > input[type=text]:nth-child(5)`).Do(ctx)
				if err != nil {
					xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}
				if t.Verbose {
					xlog.Infof("SignInTask(%v) sign start--->7", t.Username)
				}
				err = chromedp.SendKeys(`body > form > input[type=text]:nth-child(5)`, authnum).Do(ctx)
				if err != nil {
					xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}
				err = chromedp.Sleep(1 * time.Second).Do(ctx)
				if err != nil {
					xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}

				if t.Verbose {
					xlog.Infof("SignInTask(%v) sign start--->8", t.Username)
				}

				err = chromedp.EvaluateAsDevTools(`document.querySelector("input[type=submit][value='确定']").click()`, nil).Do(ctx)
				// err = chromedp.Click(`input[type=submit][value="确定"]`, chromedp.NodeVisible).Do(ctx)
				// err = chromedp.Click(`//input[@type='submit' and @value='确定']`, chromedp.NodeVisible).Do(ctx)
				// err = chromedp.Click(`body > div.mainareaOutside_pc > div.mainareaCenter_pc > form > input[type=submit]:nth-child(17)`).Do(ctx)
				if err != nil {
					if t.Verbose {
						xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					}
					err = chromedp.OuterHTML(`body`, &str).Do(ctx)
					if err != nil {
						return err
					}
				}
				if t.Verbose {
					xlog.Infof("SignInTask(%v) sign start--->9", t.Username)
				}
				outerHTML := ""
				err = XOuterHTML(ctx, &outerHTML)
				if err != nil {
					if t.Verbose {
						xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					}
					return err
				}
				if strings.Contains(outerHTML, "失败!系统检测多号刷签到!") {
					result = "失败！系统检测多号刷签到！"
					xlog.Infof("SignInTask(%v) sign with 失败!系统检测多号刷签到!", t.Username)
				} else if !strings.Contains(outerHTML, "签到成功,去记录一下") {
					result = outerHTML
					xlog.Infof("SignInTask(%v) sign with %v", t.Username, outerHTML)
				}

			default:
				err = chromedp.EvaluateAsDevTools(`document.querySelector("input[type=submit][value='确定']").click()`, nil).Do(ctx)
				// err := chromedp.Submit(`body > form > input[type=submit]:nth-child(8)`, chromedp.NodeVisible).Do(ctx)
				if err != nil {
					xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}
			}
			chromedp.Sleep(1 * time.Second).Do(ctx)
			return nil
		}),
		chromedp.Sleep(1*time.Second),
		chromedp.Navigate(`https://tx.com.cn/shop/vip/vipuse.do?gid=346`),
		chromedp.Sleep(1*time.Second),
		chromedp.Navigate(`https://tx.com.cn/in/logout.do`),
	)
	if err == nil && len(result) < 1 {
		t.failed = 0
		t.successTime = now
		xlog.Infof("SignInTask(%v) sign success", t.Username)
		conf.Conf.UpdateUser(t.Username, "signIN", time.Now().Format(`2006-01-02 15:04:05`))
		return
	}
	return
}
