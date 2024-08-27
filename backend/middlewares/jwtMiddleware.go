package middlewares

import (
	"go-react/backend/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context)  {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you are not authorised to access this resource.", "detail": "no authorization header found"})
			return
		}
		
		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

		isValid, msg := auth.ValidateJWT(bearerToken); if !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "error when validating bearer token", "detail": msg})
			return
		}

		c.Next()
	}
}