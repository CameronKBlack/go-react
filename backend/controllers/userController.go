package controllers

import (
	"context"
	"fmt"
	"go-react/backend/middlewares"
	"go-react/backend/models"
	"go-react/backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserList(client *mongo.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := middlewares.UserFromContext(c)
        if user == nil || user.Username == "" {
            c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You are not logged in or do not have access to this resource."})
            return
        }

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

func LoginLogic(username string, password string, client *mongo.Client) (string, models.User, error) {
    var result models.User
    userColl := client.Database("gore").Collection("user")

    err := userColl.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return "User not found.", models.User{}, err
        }
        return "Failed to fetch user.", result, err
    }
    if password != result.Password {
        return "Password provided does not match.", models.User{}, fmt.Errorf("incorrect password")
    }

    return "Login successful.\n", result, nil
}

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

        message, user, err := LoginLogic(credentials.Username, credentials.Password, client)
        if err != nil {
            c.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error(), "message": message})
            return
        }

        token, err := generateToken(user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message":     message,
            "userDetails": utils.ConvertUserFormat([]models.User{user}),
            "token":       token,
        })
    }
}

func generateToken(user models.User) (string, error) {
    //TODO: Implement JWT or other token generation logic
    return "valid-token", nil
}