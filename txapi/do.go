package txapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wfunc/autotx/conf"
)

func AddDoHandler(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")
	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request or Empty",
		})
		return
	}
	conf.Conf.AddDo(key, value)
	c.JSON(http.StatusOK, gin.H{
		"message": "Added",
	})
}

func RemoveDoHandler(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")
	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request or Empty",
		})
		return
	}
	conf.Conf.RemoveDo(key, value)
	c.JSON(http.StatusOK, gin.H{
		"message": "Removed",
	})
}

func ListDoHandler(c *gin.Context) {
	key := c.Query("key")
	Do := conf.Conf.ListDo(key)
	c.JSON(http.StatusOK, gin.H{
		"message": "Do",
		"do":      Do,
	})
}
