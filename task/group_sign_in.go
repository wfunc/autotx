package task

import (
	"context"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/wfunc/autotx/conf"
	"github.com/wfunc/go/xlog"
	"github.com/wfunc/util/xmap"
)

type GroupSignInTask struct {
	*BaseTask
	Name        string
	successTime time.Time
}

func NewGroupSignInTask(username, password string) *GroupSignInTask {
	base := NewBaseTaskWithUserInfo(username, password)
	base.Timeout = 120 * time.Minute
	t := &GroupSignInTask{Name: "group-sign-in", BaseTask: base}
	base.ParentTask = t
	return t
}

func (t *GroupSignInTask) Run() {
	xlog.Infof("GroupSignInTask(%v) started", t.Username)
	t.groupSignIn()
	ticker := time.NewTicker(t.TickerDelay)
	defer ticker.Stop()
	running := true
	for running {
		select {
		case <-t.exiter:
			running = false
		case <-ticker.C:
			t.groupSignIn()
		}
	}
	if t.Cancel != nil {
		t.Cancel()
	}
	xlog.Infof("GroupSignInTask(%v) done with %v", t.Username, t.Cancel)
}

func (t *GroupSignInTask) Stop() {
	t.BaseTask.stop()
}

func (t *GroupSignInTask) Info() (result xmap.M) {
	result = t.BaseInfo()
	result["success_time"] = t.successTime
	return
}

func (t *GroupSignInTask) TaskName() string {
	return t.Username + "->" + t.Name
}

func (t *GroupSignInTask) groupSignIn() (err error) {
	now := time.Now()
	makeTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 6, 0, 0, now.Location())
	if t.successTime.Year() == now.Year() && t.successTime.Month() == now.Month() && t.successTime.Day() == now.Day() {
		if t.Verbose {
			xlog.Infof("GroupSignInTask(%v) sign skipped and next will after on %v", t.Username, makeTime.AddDate(0, 0, 1).Sub(now))
		}
		if makeTime.AddDate(0, 0, 1).Sub(now) > time.Hour {
			time.Sleep(300 * time.Second)
		}
		return
	}
	// 如果now小于makeTime,直到makeTime到了再继续
	if now.Before(makeTime) {
		xlog.Infof("GroupSignInTask(%v) sign will on %v start", t.Username, makeTime)
		time.Sleep(makeTime.Sub(now))
	}
	t.CreateChromedpContext(t.Timeout)
	defer t.Cancel()
	// login
	err = t.login()
	if err != nil {
		xlog.Infof("GroupSignInTask(%v) login failed with err %v", t.Username, err)
		return
	}
	var result string
	err = chromedp.Run(t.ctx,
		chromedp.Navigate(`https://tx.com.cn/myroom/visitor/cs/summation.do?appid=1&referer=zoneIndex`),
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var outerHTML string
			err = XOuterHTML(ctx, &outerHTML)
			if err != nil {
				xlog.Infof("GroupSignInTask(%v) sign failed with err %v", t.Username, err)
				return err
			}
			if strings.Contains(outerHTML, "哎呀,今天您已领取过啦,请明天再来吧") {
				return nil
			}
			switch true {
			case strings.Contains(outerHTML, "请填写以下验证码"):
				var str string
				err = chromedp.Evaluate(`document.querySelector("body > div.mainareaOutside_pc > div.mainareaCenter_pc > img").src`, &str).Do(ctx)
				if err != nil {
					xlog.Infof("GroupSignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}

				authnum := getCode(str)
				err = chromedp.Sleep(1 * time.Second).Do(ctx)
				if err != nil {
					xlog.Infof("GroupSignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}
				if t.Verbose {
					xlog.Infof("GroupSignInTask(%v) sign start---> url %v code %v", t.Username, str, authnum)
				}
				err = chromedp.WaitVisible(`//input[@name="authnum"]`).Do(ctx)
				if err != nil {
					xlog.Infof("GroupSignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}
				if t.Verbose {
					xlog.Infof("GroupSignInTask(%v) sign start--->7", t.Username)
				}
				err = chromedp.SendKeys(`//input[@name="authnum"]`, authnum).Do(ctx)
				if err != nil {
					xlog.Infof("GroupSignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}
				err = chromedp.Sleep(1 * time.Second).Do(ctx)
				if err != nil {
					xlog.Infof("GroupSignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}

				if t.Verbose {
					xlog.Infof("GroupSignInTask(%v) sign start--->8", t.Username)
				}

				err = chromedp.EvaluateAsDevTools(`document.querySelector("input[type=submit][value='确认领取']").click()`, nil).Do(ctx)
				if err != nil {
					if t.Verbose {
						xlog.Infof("GroupSignInTask(%v) sign failed with err %v", t.Username, err)
					}
					err = chromedp.OuterHTML(`body`, &str).Do(ctx)
					if err != nil {
						return err
					}
				}
				if t.Verbose {
					xlog.Infof("GroupSignInTask(%v) sign start--->9", t.Username)
				}
				outerHTML := ""
				err = XOuterHTML(ctx, &outerHTML)
				if err != nil {
					if t.Verbose {
						xlog.Infof("GroupSignInTask(%v) sign failed with err %v", t.Username, err)
					}
					return err
				}
				if strings.Contains(outerHTML, "失败!系统检测多号刷签到!") {
					result = "失败！系统检测多号刷签到！"
					xlog.Infof("GroupSignInTask(%v) sign with 失败!系统检测多号刷签到!", t.Username)
				}
			default:
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
		t.successTime = now
		xlog.Infof("GroupSignInTask(%v) sign success", t.Username)
		conf.Conf.UpdateUser(t.Username, "groupSignIN", time.Now().Format(`2006-01-02 15:04:05`))
	}
	return

}
