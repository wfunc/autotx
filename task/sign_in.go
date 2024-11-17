package task

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type SignInTask struct {
	*BaseTask
	Name      string
	DoneAfter time.Duration
	nextAfter time.Time
	signTime  time.Time
}

func NewSignInTask(username, password string) *SignInTask {
	base := NewBaseTaskWithUserInfo(username, password)
	t := &SignInTask{Name: "sign-in", BaseTask: base}
	return t
}

func (t *SignInTask) Run() {
	t.clear()
	log.Printf("SignInTask(%v) started", t.Username)
	t.CreateChromedpContext(t.Timeout)
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
	t.Cancel()
	log.Printf("SignInTask(%v) done", t.Username)
	return
}

func (t *SignInTask) Stop() {
	t.BaseTask.stop()
}

func (t *SignInTask) TaskName() string {
	return t.Username + "->" + t.Name
}

func (t *SignInTask) clear() {
	t.signTime = time.Time{}
}

func (t *SignInTask) sign() (err error) {
	now := time.Now()
	if t.signTime.Year() == now.Year() && t.signTime.Month() == now.Month() && t.signTime.Day() == now.Day() {
		if t.Verbose {
			log.Printf("SignInTask(%v) sign skipped", t.Username)
		}
		return
	}
	// login
	err = t.login()
	if err != nil {
		log.Printf("SignInTask(%v) login failed with err %v", t.Username, err)
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
				log.Printf("SignInTask(%v) sign failed with err %v", t.Username, err)
				return err
			}
			if strings.Contains(str, "今日已签到过") {
				return nil
			}
			switch true {
			case strings.Contains(str, "请输入图片中的验证码"):
				err = chromedp.Evaluate(`document.querySelector("body > form > img").src`, &str).Do(ctx)
				if err != nil {
					log.Printf("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}

				authnum := getCode(str)
				err = chromedp.Sleep(1 * time.Second).Do(ctx)
				if err != nil {
					log.Printf("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}
				err = chromedp.WaitVisible(`body > form > input[type=text]:nth-child(5)`).Do(ctx)
				if err != nil {
					log.Printf("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}
				err = chromedp.SendKeys(`body > form > input[type=text]:nth-child(5)`, authnum).Do(ctx)
				if err != nil {
					log.Printf("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}
				err = chromedp.Sleep(1 * time.Second).Do(ctx)
				if err != nil {
					log.Printf("SignInTask(%v) sign failed with err %v", t.Username, err)
					return err
				}

				err = chromedp.Click(`body > form > input[type=submit]:nth-child(17)`).Do(ctx)
				if err != nil {
					err = chromedp.OuterHTML(`body`, &str).Do(ctx)
					if err != nil {
						return err
					}
				}

			default:
				err := chromedp.Submit(`body > div.mainareaOutside_pc > div.mainareaCenter_pc > form > input[type=submit]:nth-child(8)`, chromedp.NodeVisible).Do(ctx)
				if err != nil {
					log.Printf("SignInTask(%v) sign failed with err %v", t.Username, err)
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
		t.signTime = now
		log.Printf("SignInTask(%v) sign success", t.Username)
	}
	return
}
