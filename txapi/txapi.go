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

	task := api.Group("/task")
	task.GET("/list", ListTasksHandler)

	notDo := api.Group("notDo")
	notDo.GET("/add", AddNotDoHandler)
	notDo.GET("/remove", RemoveNotDoHandler)
	notDo.GET("/list", ListNotDoHandler)

	do := api.Group("do")
	do.GET("/add", AddDoHandler)
	do.GET("/remove", RemoveDoHandler)
	do.GET("/list", ListDoHandler)
}
