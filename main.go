package main

import (
	"os"
	"thoughts-api/src/middleware"
	"thoughts-api/src/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)

	router.Use(middleware.Authentication())
	routes.ThoughtsRoutes(router)
	routes.RepliesRoutes(router)
	router.Run(":" + port)
}
