package controllers

import (
	"context"
	"go-react/backend/models"
	"go-react/backend/services"
	"go-react/backend/states"
	"go-react/backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoginCall(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&credentials); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		message, u, err := services.LoginLogic(credentials.Username, credentials.Password, client)
		if err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err, "user": u, "message": message})
			return
		}

		states.UserMutex.Lock()
		states.CurrentUser = u
		states.UserMutex.Unlock()
		userResponse := utils.ConvertUserFormat([]models.User{u})
		c.String(http.StatusOK, message)
		c.IndentedJSON(http.StatusOK, gin.H{"userDetails": userResponse})
	}
}

func GetUserList(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userColl := client.Database("gore").Collection("user")
		var users []bson.M
		cursor, err := userColl.Find(context.TODO(), bson.D{})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch users."})
			return
		}
		defer cursor.Close(context.TODO())

		if err = cursor.All(context.TODO(), &users); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to decode users"})
			return
		}

		userList, err := utils.ConvertFromBSONToUserSlice(users)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to convert from BSON to user list"})
			return
		}
		userResponses := utils.ConvertUserFormat(userList)
		c.IndentedJSON(http.StatusOK, userResponses)
	}
}
