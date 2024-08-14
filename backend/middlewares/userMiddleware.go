package middlewares

import (
	"go-react/backend/states"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		states.UserMutex.RLock()
		defer states.UserMutex.RUnlock()

		if states.CurrentUser.Username == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You are not authorised to access this resource."})
			return
		}

		c.Set("currentUser", states.CurrentUser)
		c.Next()
	}
}