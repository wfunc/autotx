package task

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

// BaseTask 包含公共的任务配置和行为
type BaseTask struct {
	Username    string
	Password    string
	Headless    bool
	UserAgent   string
	Proxy       string
	ctx         context.Context
	Cancel      context.CancelFunc
	started     bool
	Timeout     time.Duration
	TickerDelay time.Duration
	exiter      chan int
	Verbose     bool
}

func NewBaseTask(ctx context.Context) *BaseTask {
	return &BaseTask{
		ctx:         ctx,
		Headless:    true,
		UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
		Timeout:     30 * time.Second,
		TickerDelay: 5 * time.Second,
		exiter:      make(chan int, 2),
	}
}

func NewBaseTaskWithUserInfo(username, password string) *BaseTask {
	ctx, cancel := context.WithCancel(context.Background())
	b := &BaseTask{
		ctx:         ctx,
		Cancel:      cancel,
		Username:    username,
		Password:    password,
		Headless:    true,
		UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
		Timeout:     30 * time.Second,
		TickerDelay: 5 * time.Second,
		exiter:      make(chan int, 2),
	}
	return b
}

func (b *BaseTask) CreateChromedpContext(timeout time.Duration) {
	if b.started {
		return
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", b.Headless),
		chromedp.Flag("window-size", "400,745"),
		chromedp.UserAgent(b.UserAgent),
	)
	if len(b.Proxy) > 0 {
		opts = append(opts, chromedp.ProxyServer(b.Proxy), chromedp.Flag("proxy-bypass-list", "<-loopback>"))
	}
	if b.Headless {
		opts = append(opts, chromedp.DisableGPU)
	}
	allocCtx, cancel := chromedp.NewExecAllocator(b.ctx, opts...)
	taskCtx, cancelCtx := chromedp.NewContext(allocCtx)
	// taskCtx, cancelTimeout := context.WithTimeout(taskCtx, timeout) // timeout
	b.ctx = taskCtx
	b.Cancel = func() {
		// cancelTimeout()
		cancelCtx()
		cancel()
		b.started = false
	}
	b.started = true
}

func (b *BaseTask) stop() {
	b.exiter <- 0
	b.started = false
}
