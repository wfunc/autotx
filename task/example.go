package task

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/wfunc/go/xlog"
)

type ExampleTask struct {
	*BaseTask
	Name      string
	DoneAfter time.Duration
}

func NewExampleTask(name string) *ExampleTask {
	ctx, cancel := context.WithCancel(context.Background())
	base := NewBaseTask(ctx)
	base.ctx = ctx
	base.Cancel = cancel
	t := &ExampleTask{Name: name, BaseTask: base}
	return t
}

func (t *ExampleTask) Run() {
	xlog.Infof("Task %s started", t.Name)

	// 创建 Chrome 执行上下文
	t.CreateChromedpContext(30 * time.Second)
	defer t.Cancel()

	for i := 0; i < 3; i++ {
		// 执行任务逻辑
		var title string
		var err = chromedp.Run(t.ctx,
			chromedp.Navigate("https://example.com"),
			chromedp.Title(&title),
		)
		if err != nil {
			xlog.Infof("Task %s error: %v", t.Name, err)
			return
		}
		xlog.Infof("Page title: %s", title)
		time.Sleep(1 * time.Second)
	}

	// 延迟结束
	if t.DoneAfter > 0 {
		time.Sleep(t.DoneAfter)
	}
	return
}

func (t *ExampleTask) Stop() {
	t.BaseTask.stop()
}

func (t *ExampleTask) TaskName() string {
	return t.Username + "->" + t.Name
}
