package services

import (
	"bytes"
	"context"
	"fmt"
	"go-react/backend/models"
	"go-react/backend/utils"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterNewUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
        body, err := io.ReadAll(c.Request.Body)
        if err != nil {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
            return
        }

		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		requiredFields := utils.GetRequiredFields(reflect.TypeOf(models.User{}))

		if !utils.CheckRequiredFields(body, requiredFields) {
			return
		}
		
		var newUser models.User	

		if err := c.BindJSON(&newUser); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		coll := client.Database("gore").Collection("user")

		var existingUser models.User
		err = coll.FindOne(context.TODO(), bson.M{"username": newUser.Username}).Decode(&existingUser)
	
		if err == mongo.ErrNoDocuments {
			res, err := coll.InsertOne(context.TODO(), newUser)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
				return
			}
			c.IndentedJSON(http.StatusOK, gin.H{"inserted_id": res.InsertedID})
		} else if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred.", "details": err.Error()})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "A user with these details already exists on our system."})
		}
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

	return "Login successful.\n\n", result, nil
}