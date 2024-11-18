package txapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wfunc/autotx/conf"
	"github.com/wfunc/autotx/runner"
)

func ListUsersHandler(c *gin.Context) {
	users := conf.Conf.GetUsers()
	seeds := conf.Conf.GetSeedsRevert()
	for u, uc := range users {
		setSeeds := uc.Map("set_seeds")
		for k := range setSeeds {
			setSeeds[k] = seeds[k]
		}
		uc["set_seeds"] = setSeeds
		users[u] = uc
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Users",
		"users":   users,
	})
}

func AddUserHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request or Empty",
		})
		return
	}
	conf.Conf.AddUser(username, password)
	runner.RunnerShared.Reload(username)
	c.JSON(http.StatusOK, gin.H{
		"message": "Added",
	})
}

func RemoveUserHandler(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request or Empty",
		})
		return
	}
	conf.Conf.RemoveUser(username)
	runner.RunnerShared.StopTask(username)
	c.JSON(http.StatusOK, gin.H{
		"message": "Removed",
	})
}

func ReloadUserHandler(c *gin.Context) {
	username := c.Query("username")
	if len(username) > 0 {
		runner.RunnerShared.StopTask(username)
		runner.RunnerShared.Reload(username)
	} else {
		users := conf.Conf.GetUsers()
		for username := range users {
			runner.RunnerShared.StopTask(username)
			runner.RunnerShared.Reload(username)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Reloaded",
	})
}
