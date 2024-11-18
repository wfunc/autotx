package txapi

import (
	"github.com/gin-gonic/gin"
	"github.com/wfunc/autotx/conf"
	"github.com/wfunc/autotx/runner"
)

func Handle(router *gin.Engine) {
	api := router.Group("/api")
	user := api.Group("/user")
	user.GET("/list", UsersHandler)
	user.GET("/add", AddUserHandler)
	user.GET("/remove", RemoveUserHandler)
}

func UsersHandler(c *gin.Context) {
	users := conf.Conf.GetUsers()
	c.JSON(200, gin.H{
		"message": "Users",
		"users":   users,
	})
}

func AddUserHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || password == "" {
		c.JSON(400, gin.H{
			"message": "Bad Request or Empty",
		})
		return
	}
	conf.Conf.AddUser(username, password)
	runner.RunnerShared.Reload(username)
	c.JSON(200, gin.H{
		"message": "Added",
	})
}

func RemoveUserHandler(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(400, gin.H{
			"message": "Bad Request or Empty",
		})
		return
	}
	conf.Conf.RemoveUser(username)
	runner.RunnerShared.StopTask(username)
	c.JSON(200, gin.H{
		"message": "Removed",
	})
}
