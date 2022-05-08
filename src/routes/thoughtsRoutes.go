package routes

import (
	"thoughts-api/src/controllers"

	"github.com/gin-gonic/gin"
)

func ThoughtsRoutes(router *gin.Engine) {
	router.POST("/thoughts/add", controllers.AddThought())
	router.GET("/thoughts/mine", controllers.ListMyThoughts())
	router.GET("/thoughts/:username", controllers.ListOtherUserThoughts())
	router.DELETE("/thoughts/delete/:id", controllers.DeleteThoughts())
}
