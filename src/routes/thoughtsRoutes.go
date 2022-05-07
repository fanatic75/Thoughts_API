package routes

import (
	"thoughts-api/src/controllers"

	"github.com/gin-gonic/gin"
)

func ThoughtsRoutes(route *gin.Engine) {
	route.POST("/thoughts", controllers.AddThought())
}
