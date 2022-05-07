package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	helper "thoughts-api/src/helpers"
)

// Auth validates token and authorizes users
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")

		if bearerToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}
		clientToken := strings.Split(bearerToken, " ")
		if len(clientToken) < 2 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Token Provided")})
			c.Abort()
			return
		}
		claims, err := helper.ValidateToken(clientToken[1])
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token not valid"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)

		c.Next()

	}
}
