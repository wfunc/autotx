package txapi

import (
	"github.com/gin-gonic/gin"
)

func Handle(router *gin.Engine) {
	api := router.Group("/api")
	user := api.Group("/user")
	user.GET("/list", ListUsersHandler)
	user.GET("/add", AddUserHandler)
	user.GET("/remove", RemoveUserHandler)
	user.GET("/reload", ReloadUserHandler)
}
