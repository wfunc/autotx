package task

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type ExchangeCardTask struct {
	*BaseTask
	Name        string
	DoneAfter   time.Duration
	successTime time.Time
}

func NewExchangeCardTask(username, password string) *ExchangeCardTask {
	base := NewBaseTaskWithUserInfo(username, password)
	t := &ExchangeCardTask{Name: "exchange-card", BaseTask: base}
	return t
}

func (t *ExchangeCardTask) Run() {
	log.Printf("ExchangeCardTask(%v) started", t.Username)
	t.CreateChromedpContext(t.Timeout)
	t.exchange()
	ticker := time.NewTicker(t.TickerDelay)
	defer ticker.Stop()
	running := true
	for running {
		select {
		case <-t.exiter:
			running = false
		case <-ticker.C:
			t.exchange()
		}
	}
	t.Cancel()
	log.Printf("ExchangeCardTask(%v) done", t.Username)
}

func (t *ExchangeCardTask) Stop() {
	t.BaseTask.stop()
}

func (t *ExchangeCardTask) TaskName() string {
	return t.Username + "->" + t.Name
}

func (t *ExchangeCardTask) exchange() (err error) {
	now := time.Now()
	if t.successTime.Year() == now.Year() && t.successTime.Month() == now.Month() && t.successTime.Day() == now.Day() {
		if t.Verbose {
			log.Printf("SignInTask(%v) sign skipped", t.Username)
		}
		return
	}
	// login
	err = t.login()
	if err != nil {
		log.Printf("ExchangeCardTask(%v) login failed with err %v", t.Username, err)
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
					chromedp.Navigate(url).Do(ctx)
					time.Sleep(1 * time.Second)
					var outHTML string
					err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outHTML).Do(ctx)
					if err != nil {
						return err
					}
					if strings.Contains(outHTML, "幸运卡:") {
						err = chromedp.Click(fmt.Sprintf(`body > div.mainareaOutside_pc > div.mainareaCenter_pc > div.mainarea > div > table > tbody > tr > td:nth-child(%s) > form > input[type=submit]:nth-child(2)`, subCard)).Do(ctx)
						if err != nil {
							return err
						}
						err = chromedp.OuterHTML(`body > div.mainareaOutside_pc > div.mainareaCenter_pc`, &outHTML).Do(ctx)
						if err != nil {
							return err
						}
						if strings.Contains(outHTML, "恭喜") {
							if t.Verbose {
								log.Printf("ExchangeCardTask(%v) Success", t.Username)
							}
							break
						} else if strings.Contains(outHTML, "厉害!") {
							if t.Verbose {
								log.Printf("ExchangeCardTask(%v) All Success done", t.Username)
							}
							return nil
						} else {
							if t.Verbose {
								log.Printf("ExchangeCardTask(%v) Failed", t.Username)
							}
						}
					} else {
						if t.Verbose {
							log.Printf("ExchangeCardTask(%v) No luck", t.Username)
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
	}
	return
}
