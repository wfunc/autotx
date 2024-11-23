package task

import (
	"context"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/wfunc/go/xlog"
)

func (t *FarmTask) water() (err error) {
	count := 0
	f := chromedp.ActionFunc(func(ctx context.Context) error {
		waterUrl := "https://tx.com.cn/plugins/farm/rank.do?oper=3"
		err = chromedp.Navigate(waterUrl).Do(ctx)
		if err != nil {
			return err
		}
		tasks := NewM()
		var outerHTML string
		for {
			err = XOuterHTML(ctx, &outerHTML)
			if err != nil {
				if t.Verbose {
					xlog.Infof("FarmTask(%v-%v) Failed with err %v", t.Username, t.target, err)
				}
				return err
			}
			tasks.ExtractLinksWithPrefix(outerHTML, "https://tx.com.cn", []string{"/plugins/farm/index.do?uid="})
			if t.Verbose {
				xlog.Infof("FarmTask(%v-%v) Success tasks ---> %v", t.Username, t.target, len(tasks.Map))
			}
			if !strings.Contains(outerHTML, "下页") {
				break
			}
			selector := `//a[contains(text(), "下页")]`
			if !t.clickNext(ctx, selector) {
				break
			}
		}
		if t.Verbose {
			for count, m := range tasks.CountMap {
				xlog.Infof("count %v m %v", count, len(m))
			}
			xlog.Infof("FarmTask(%v-%v) Success tasks ---> %v", t.Username, t.target, len(tasks.Map))
		}
		for k, v := range tasks.Map {
			err = chromedp.Navigate(v + k).Do(ctx)
			if err != nil {
				if t.Verbose {
					xlog.Infof("FarmTask(%v-%v) Navigate(%v) Failed with err %v", t.Username, t.target, v+k, err)
				}
				return err
			}
			count++
			if t.Verbose {
				xlog.Infof("FarmTask(%v-%v) Success ---> %v", t.Username, t.target, v+k)
			}
			subTasks := NewM()
			km := map[string]bool{}
			for {
				err = XOuterHTML(ctx, &outerHTML)
				if err != nil {
					return err
				}
				km1, _ := subTasks.ExtractLinksWithPrefix(outerHTML, "https://tx.com.cn/plugins/farm/", []string{"water.do", "killInsects.do", "weeding.do"})
				for kmk, kmv := range km1 {
					km[kmk] = kmv
				}
				if !strings.Contains(outerHTML, "下页") {
					// if t.Verbose {
					// 	xlog.Infof("FarmTask(%v-%v) Success outerHTML ---> %v", t.Username, t.target, outerHTML)
					// }
					break
				}
				selector := `//a[contains(text(), "下页")]`
				if !t.clickNext(ctx, selector) {
					break
				}
			}
			waterDone := false
			if t.Verbose {
				xlog.Infof("FarmTask(%v-%v) Success subTasks ---> %v", t.Username, t.target, len(subTasks.Map))
			}
			for subK, subV := range subTasks.Map {
				if strings.Contains(subK, "water.do") && km["water.go"] && waterDone {
					continue
				}
				err = chromedp.Navigate(subV + subK).Do(ctx)
				if err != nil {
					if t.Verbose {
						xlog.Infof("FarmTask(%v-%v) Failed with err %v", t.Username, t.target, err)
					}
					return err
				}
				if t.Verbose {
					xlog.Infof("FarmTask(%v-%v) Success subTask ---> %v", t.Username, t.target, v+k)
				}
				if strings.Contains(subK, "water.do") && km["water.do"] && !waterDone {
					err = XOuterHTML(ctx, &outerHTML)
					if err != nil {
						return err
					}
					switch true {
					case strings.Contains(outerHTML, "需要浇水"):
						// for {
						// 	nexts := NewM()

						// 	nexts.ExtractLinksWithPrefix(outerHTML, "https://tx.com.cn/plugins/farm/", []string{"water.do"})

						// 	if len(nexts.Map) < 1 {
						// 		break
						// 	}
						// 	for nextK, nextV := range nexts.Map {
						// 		err = chromedp.Navigate(nextV + nextK).Do(ctx)
						// 		if err != nil {
						// 			return err
						// 		}
						// 		if t.Verbose {
						// 			xlog.Infof("FarmTask(%v-%v) Success subTask next ---> %v", t.Username, t.target, v+k)
						// 		}
						// 		err = XOuterHTML(ctx, &outerHTML)
						// 		if err != nil {
						// 			return err
						// 		}
						// 		if !strings.Contains(outerHTML, "需要浇水") {
						// 			break
						// 		}
						// 	}
						// }
					default:
						// waterDone = true
					}
				}

			}
		}
		return nil
	})
	err = chromedp.Run(t.ctx, f)
	if err != nil {
		if t.Verbose {
			xlog.Infof("FarmTask(%v-%v) Failed with err %v", t.Username, t.target, err)
		}
		return err
	}
	if t.Verbose {
		xlog.Infof("FarmTask(%v-%v) Success with %v", t.Username, t.target, count)
	}
	return
}
