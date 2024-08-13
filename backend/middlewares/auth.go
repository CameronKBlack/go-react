package middlewares

import (
	"context"
	"go-react/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

const userKey = "user"


func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "No authorization token provided"})
            c.Abort()
            return
        }

        user, err := authenticateUser(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
            c.Abort()
            return
        }

        ctx := context.WithValue(c.Request.Context(), userKey, user)
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}


func authenticateUser(token string) (*models.User, error) {
    if token == "valid-token" {
        return &models.User{
            Username: "exampleUser",
            // Other user fields...
        }, nil
    }
    return nil, nil
}

func UserFromContext(c *gin.Context) *models.User {
    return c.Request.Context().Value(userKey).(*models.User)
}
