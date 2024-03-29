package routes

import (
	controller "thoughts-api/src/controllers"

	"github.com/gin-gonic/gin"
)

// UserRoutes
func UserRoutes(route *gin.Engine) {
	route.POST("/users/signup", controller.SignUp())
	route.POST("/users/login", controller.Login())
}
