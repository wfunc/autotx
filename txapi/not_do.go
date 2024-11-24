package txapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wfunc/autotx/conf"
)

func AddNotDoHandler(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")
	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request or Empty",
		})
		return
	}
	conf.Conf.AddNotDo(key, value)
	c.JSON(http.StatusOK, gin.H{
		"message": "Added",
	})
}

func RemoveNotDoHandler(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")
	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request or Empty",
		})
		return
	}
	conf.Conf.RemoveNotDo(key, value)
	c.JSON(http.StatusOK, gin.H{
		"message": "Removed",
	})
}

func ListNotDoHandler(c *gin.Context) {
	key := c.Query("key")
	notDo := conf.Conf.ListNotDo(key)
	c.JSON(http.StatusOK, gin.H{
		"message": "NotDo",
		"not_do":  notDo,
	})
}
