package task

import (
	"context"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/wfunc/go/xlog"
	"github.com/wfunc/util/xmap"
)

type SignInTask struct {
	*BaseTask
	Name        string
	DoneAfter   time.Duration
	successTime time.Time
}

func NewSignInTask(username, password string) *SignInTask {
	base := NewBaseTaskWithUserInfo(username, password)
	base.Timeout = 5 * time.Minute
	t := &SignInTask{Name: "sign-in", BaseTask: base}
	return t
}

func (t *SignInTask) Run() {
	t.clear()
	xlog.Infof("SignInTask(%v) started", t.Username)

	t.sign()
	ticker := time.NewTicker(t.TickerDelay)
	defer ticker.Stop()
	running := true
	for running {
		select {
		case <-t.exiter:
			running = false
		case <-ticker.C:
			t.sign()
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
}

func (t *SignInTask) sign() (err error) {
	now := time.Now()
	if t.successTime.Year() == now.Year() && t.successTime.Month() == now.Month() && t.successTime.Day() == now.Day() {
		if t.Verbose {
			xlog.Infof("SignInTask(%v) sign skipped", t.Username)
		}
		return
	}
	t.CreateChromedpContext(t.Timeout)
	defer t.Cancel()
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
		chromedp.Navigate(`https://tx.com.cn/activity/qq/cs/sign.do`),
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var str string
			err = chromedp.OuterHTML(`body`, &str).Do(ctx)
			if err != nil {
				xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
				return err
			}
			if strings.Contains(str, "今日已签到过") {
				return nil
			}
			switch true {
			case strings.Contains(str, "请输入图片中的验证码"):
				err = chromedp.Evaluate(`document.querySelector("body > div.mainareaOutside_pc > div.mainareaCenter_pc > form > img").src`, &str).Do(ctx)
				if err != nil {
					xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}

				authnum := getCode(str)
				err = chromedp.Sleep(1 * time.Second).Do(ctx)
				if err != nil {
					xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}

				err = chromedp.WaitVisible(`body > div.mainareaOutside_pc > div.mainareaCenter_pc > form > input[type=text]:nth-child(5)`).Do(ctx)
				if err != nil {
					xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}
				err = chromedp.SendKeys(`body > div.mainareaOutside_pc > div.mainareaCenter_pc > form > input[type=text]:nth-child(5)`, authnum).Do(ctx)
				if err != nil {
					xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}
				err = chromedp.Sleep(1 * time.Second).Do(ctx)
				if err != nil {
					xlog.Infof("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}

				err = chromedp.Click(`body > div.mainareaOutside_pc > div.mainareaCenter_pc > form > input[type=submit]:nth-child(17)`).Do(ctx)
				if err != nil {
					err = chromedp.OuterHTML(`body`, &str).Do(ctx)
					if err != nil {
						return err
					}
				}

			default:
				err := chromedp.Submit(`body > div.mainareaOutside_pc > div.mainareaCenter_pc > form > input[type=submit]:nth-child(8)`, chromedp.NodeVisible).Do(ctx)
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
	if err == nil {
		t.successTime = now
		xlog.Infof("SignInTask(%v) sign success", t.Username)
	}
	return
}
