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

type ExchangeCardTask struct {
	*BaseTask
	Name        string
	DoneAfter   time.Duration
	successTime time.Time
}

func NewExchangeCardTask(username, password string) *ExchangeCardTask {
	base := NewBaseTaskWithUserInfo(username, password)
	base.Timeout = 120 * time.Minute
	t := &ExchangeCardTask{Name: "exchange-card", BaseTask: base}
	base.ParentTask = t
	return t
}

func (t *ExchangeCardTask) Run() {
	xlog.Infof("ExchangeCardTask(%v) started", t.Username)
	t.exchangeCard()
	ticker := time.NewTicker(t.TickerDelay)
	defer ticker.Stop()
	running := true
	for running {
		select {
		case <-t.exiter:
			running = false
		case <-ticker.C:
			t.exchangeCard()
		}
	}
	if t.Cancel != nil {
		t.Cancel()
	}
	xlog.Infof("ExchangeCardTask(%v) done", t.Username)
}

func (t *ExchangeCardTask) Stop() {
	t.BaseTask.stop()
}

func (t *ExchangeCardTask) Info() (result xmap.M) {
	result = t.BaseInfo()
	result["success_time"] = t.successTime
	return
}

func (t *ExchangeCardTask) TaskName() string {
	return t.Username + "->" + t.Name
}

func (t *ExchangeCardTask) exchangeCard() (err error) {
	now := time.Now()
	if t.successTime.Year() == now.Year() && t.successTime.Month() == now.Month() && t.successTime.Day() == now.Day() {
		if t.Verbose {
			xlog.Infof("ExchangeCardTask(%v) skipped", t.Username)
		}
		return
	}
	if !conf.Conf.IsAllSignIN() {
		xlog.Infof("ExchangeCardTask(%v) All not sign in will check next 1 minute", t.Username)
		time.Sleep(1 * time.Minute)
		return
	}
	lock := ChromeManagerInstance.GetUserLock(t.Username)
	lock.Lock()
	defer lock.Unlock()
	t.CreateChromedpContext(t.Timeout)
	defer t.Cancel()
	// login
	err = t.login()
	if err != nil {
		xlog.Infof("ExchangeCardTask(%v) login failed with err %v", t.Username, err)
		return
	}

	err = chromedp.Run(t.ctx,
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for {
				url := ""
				cards := []string{"3", "5", "7", "9", "11"}
				for _, subCard := range cards {
					url = "https://tx.com.cn/room/rindex.do?op=2&ar1=696"
					err = chromedp.Navigate(url).Do(ctx)
					if err != nil {
						if t.Verbose {
							xlog.Infof("ExchangeCardTask(%v) ExchangeCardTask(%v) Failed with err %v", t.Username, subCard, err)
						}
						return err
					}
					time.Sleep(1 * time.Second)
					var outerHTML string
					err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
					if err != nil {
						if t.Verbose {
							xlog.Infof("ExchangeCardTask(%v) ExchangeCardTask(%v) Failed with err %v", t.Username, subCard, err)
						}
						return err
					}
					if strings.Contains(outerHTML, "幸运卡:") {
						err = chromedp.Click(fmt.Sprintf(`body > div.mainareaOutside_pc > div.mainareaCenter_pc > div.mainarea > div > table > tbody > tr > td:nth-child(%s) > form > input[type=submit]:nth-child(2)`, subCard)).Do(ctx)
						if err != nil {
							return err
						}
						err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outerHTML).Do(ctx)
						if err != nil {
							return err
						}
						if strings.Contains(outerHTML, "恭喜") {
							if t.Verbose {
								xlog.Infof("ExchangeCardTask(%v) Success", t.Username)
							}
							break
						} else if strings.Contains(outerHTML, "厉害!") {
							if t.Verbose {
								xlog.Infof("ExchangeCardTask(%v) All Success done", t.Username)
							}
							return nil
						} else {
							if t.Verbose {
								xlog.Infof("ExchangeCardTask(%v) Failed", t.Username)
							}
						}
					} else {
						if t.Verbose {
							xlog.Infof("ExchangeCardTask(%v) No luck", t.Username)
						}
						err = chromedp.Evaluate(`document.querySelector('body > div.mainareaOutside_pc > div.mainareaCenter_pc > div.mainarea > div > div.news > form > a:nth-child(8)').href`, &url).Do(ctx)
						if err != nil {
							return err
						}
					}
				}
				time.Sleep(1 * time.Second)
			}
		}),
	)
	if err == nil {
		t.successTime = time.Now()
		conf.Conf.UpdateUser(t.Username, "exchangeCard", t.successTime.Format(`2006-01-02 15:04:05`))
		xlog.Infof("ExchangeCardTask(%v) success with %v", t.Username, t.successTime)
	}
	return
}
