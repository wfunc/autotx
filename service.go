package main

import (
	"log"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/wfunc/autotx/conf"
	"github.com/wfunc/autotx/runner"
	"github.com/wfunc/autotx/txapi"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	authorized := gin.BasicAuth(gin.Accounts{
		"dev": "123",
	})
	r.Use(authorized)
	pprofGroup := r.Group("/")
	pprofGroup.Use(authorized)
	pprof.Register(pprofGroup)
	// api handle
	txapi.Handle(r)
	r.GET("/stop", func(c *gin.Context) {
		all := runner.RunnerShared.Stop()
		c.JSON(200, gin.H{
			"message": "Stopped",
			"all":     all,
		})
	})
	r.GET("/start", func(c *gin.Context) {
		runner.Bootstrap()
		c.JSON(200, gin.H{
			"message": "Started",
		})
	})
	// conf
	conf.Bootstrap()
	// runner
	runner.Bootstrap()
	log.Printf("Server started on :8080")
	r.Run()
}
