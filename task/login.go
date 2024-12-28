package task

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/wfunc/go/xlog"
	"github.com/wfunc/util/xhttp"
)

func (b *BaseTask) login() (err error) {
	var res string
	err = chromedp.Run(b.ctx,
		mobileSubmit("https://tx.com.cn/in/cs/login.do?type=passwd", `//input[@name="useruid"]`, b.Username, `//input[@name="password"]`, b.Password, &res),
	)
	if err != nil {
		xlog.Infof("Login(%v) failed with err %v", b.Username, err)
		return
	}
	if strings.Contains(res, "继续访问空间") {
		xlog.Infof("Login(%v) success", b.ParentTask.TaskName())
		return nil
	}
	return
}

func mobileSubmit(urlstr, sel, q, sel2, q2 string, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(sel),
		chromedp.Clear(`//input[@name="useruid"]`),
		chromedp.SendKeys(sel, q),
		chromedp.WaitVisible(sel2),
		chromedp.SendKeys(sel2, q2),
		chromedp.Submit(sel),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(1 * time.Second)
			err := chromedp.OuterHTML(`body`, res).Do(ctx)
			if err != nil {
				xlog.Infof("mobileSubmit(%s) failed with err %v", q, err)
				return err
			}
			if strings.Contains(*res, "继续访问空间") {
				return nil
			}
			err = chromedp.Evaluate(`document.querySelector('body > div:nth-child(2) > div > form:nth-child(1) > img').src`, res).Do(ctx)
			if err != nil {
				return err
			}
			authnum := getCode(*res)
			err = chromedp.SendKeys(`body > div:nth-child(2) > div > form:nth-child(1) > input[type=text]:nth-child(13)`, authnum).Do(ctx)
			if err != nil {
				return err
			}
			// ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
			// defer cancel()

			// err = chromedp.WaitVisible(`body > div:nth-child(2) > div > form:nth-child(1) > input[type=submit]:nth-child(17)`, chromedp.ByQuery).Do(ctx)
			// if err != nil {
			// 	xlog.Errorf("Failed to wait for submit button: %v", err)
			// 	return err
			// }
			// err = chromedp.Submit(`body > div:nth-child(2) > div > form:nth-child(1) > input[type=submit]:nth-child(17)`).Do(ctx)
			// xlog.Infof("mobileSubmit(%s) authnum is %s err %v", q, authnum, err)
			err = chromedp.OuterHTML(`body`, res).Do(ctx)
			if err != nil {
				xlog.Infof("mobileSubmit(%s) failed with err %v", q, err)
				return err
			}

			if strings.Contains(*res, "继续访问空间") {
				return nil
			}
			return err
		}),
	}
}

func getCode(url string) string {
	code := getCodeLogic(url)
	if len(code) != 4 {
		code = getCodeLogic(url)
	}
	return code
}

func getCodeLogic(url string) string {
	if CodeAPI {
		data, _ := xhttp.GetText(CodeURL+"%v", url)
		return data
	} else {
		// old code will be removed in the future
		args := []string{"../ddddocr/ocr.py", url}
		cmd := exec.Command("python3", args...)
		cmdData, err := cmd.CombinedOutput()
		if err != nil {
			xlog.Infof("python err is %v", err)
			return ""
		}
		return string(cmdData)
	}

}
