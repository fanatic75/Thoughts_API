package routes

import (
	"thoughts-api/src/controllers"

	"github.com/gin-gonic/gin"
)

func RepliesRoutes(router *gin.Engine) {
	router.POST("/replies/:thoughtID/add", controllers.AddReplies())
	router.DELETE("/replies/:thoughtID/delete/:id", controllers.DeleteReplies())
}
