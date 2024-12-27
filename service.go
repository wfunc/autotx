package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/wfunc/autotx/conf"
	"github.com/wfunc/autotx/runner"
	"github.com/wfunc/autotx/task"
	"github.com/wfunc/autotx/txapi"
	"github.com/wfunc/go/xlog"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.ForwardedByClientIP = true
	rate := limiter.Rate{Period: time.Minute, Limit: 100} // 100 requests per minute by IP
	store := memory.NewStore()
	middleware := mgin.NewMiddleware(limiter.New(store, rate))
	r.Use(middleware)
	// set basic auth
	authorized := gin.BasicAuth(gin.Accounts{
		"dev": "123",
	})
	r.Use(authorized)
	// set pprof
	pprof.Register(r)
	// api handle
	txapi.Handle(r)
	r.GET("/stop", func(c *gin.Context) {
		all := runner.RunnerShared.Stop()
		c.JSON(http.StatusOK, gin.H{
			"message": "Stopped",
			"all":     all,
		})
	})
	r.GET("/start", func(c *gin.Context) {
		runner.Bootstrap()
		c.JSON(http.StatusOK, gin.H{
			"message": "Started",
		})
	})
	r.GET("/robots.txt", func(c *gin.Context) {
		c.String(http.StatusNotFound, "")
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusNotFound, "")
	})
	task.BootstrapChromeManagerInstance()

	task.BootstrapConfig()
	// conf
	conf.Bootstrap()
	// runner
	runner.Bootstrap()
	xlog.Infof("Server started on :8080")
	r.Run()
}
