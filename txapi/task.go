package txapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wfunc/autotx/runner"
)

func ListTasksHandler(c *gin.Context) {
	tasks := runner.RunnerShared.ListTask()
	c.JSON(http.StatusOK, gin.H{
		"message": "Tasks",
		"tasks":   tasks,
	})
}
