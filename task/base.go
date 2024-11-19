package task

import (
	"context"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/wfunc/go/xlog"
	"github.com/wfunc/util/xmap"
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
	lock        sync.RWMutex
}

func NewBaseTask(ctx context.Context) *BaseTask {
	return &BaseTask{
		ctx:         ctx,
		Headless:    true,
		UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
		Timeout:     30 * time.Second,
		TickerDelay: 5 * time.Second,
		exiter:      make(chan int, 3),
		lock:        sync.RWMutex{},
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
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.started {
		xlog.Infof("BaseTask(%v) already started", b.Username)
		return
	}
	// opts := append(chromedp.DefaultExecAllocatorOptions[:],
	// 	chromedp.Flag("headless", b.Headless),
	// 	chromedp.Flag("window-size", "400,745"),
	// 	chromedp.UserAgent(b.UserAgent),
	// )
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", b.Headless),                      // 必须运行在无头模式
		chromedp.Flag("disable-gpu", true),                         // 禁用 GPU，加速无头模式
		chromedp.Flag("no-sandbox", true),                          // 禁用沙箱，降低资源占用
		chromedp.Flag("disable-extensions", true),                  // 禁用扩展
		chromedp.Flag("disable-images", true),                      // 禁用图片加载
		chromedp.Flag("disable-default-apps", true),                // 禁用默认应用
		chromedp.Flag("disable-dev-shm-usage", true),               // 禁用 /dev/shm 共享内存限制
		chromedp.Flag("disable-setuid-sandbox", true),              // 禁用 setuid 沙箱
		chromedp.Flag("disable-infobars", true),                    // 禁用信息栏
		chromedp.Flag("disable-popup-blocking", true),              // 禁用弹窗拦截
		chromedp.Flag("disable-translate", true),                   // 禁用翻译功能
		chromedp.Flag("disable-background-timer-throttling", true), // 禁用后台计时器限制
		chromedp.Flag("disable-renderer-backgrounding", true),      // 禁用后台渲染
		chromedp.Flag("disable-background-networking", true),       // 禁用后台网络请求
		chromedp.Flag("disable-sync", true),                        // 禁用同步
		chromedp.Flag("mute-audio", true),                          // 禁用音频
		// chromedp.Flag("remote-debugging-port", "9222"),             // 开启远程调试端口（可选）
		chromedp.Flag("window-size", "400,745"), // 设置窗口大小
		chromedp.UserAgent(b.UserAgent),         // 设置用户代理
	)
	if len(b.Proxy) > 0 {
		opts = append(opts, chromedp.ProxyServer(b.Proxy), chromedp.Flag("proxy-bypass-list", "<-loopback>"))
	}
	if b.Headless {
		opts = append(opts, chromedp.DisableGPU)
	}
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	taskCtx, cancelCtx := chromedp.NewContext(allocCtx)
	withTimeoutCtx, cancelTimeout := context.WithTimeout(taskCtx, timeout) // timeout
	b.ctx = withTimeoutCtx
	b.Cancel = func() {
		cancelTimeout()
		cancelCtx()
		cancel()
		b.lock.Lock()
		b.started = false
		b.Cancel = nil
		b.lock.Unlock()
	}
	b.started = true
}

func (b *BaseTask) stop() {
	b.exiter <- 0
}

func (b *BaseTask) BaseInfo() xmap.M {
	b.lock.RLock()
	defer b.lock.RUnlock()
	result := xmap.M{}
	result["started"] = b.started
	return result
}
