package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-react/backend/models"
	"go-react/backend/utils"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

func HashAndSalt(password string) string {
	salt := make([]byte, 16)

	hash := argon2.IDKey([]byte(password), salt, 1, (32*1024), 1, 32)
	
	return base64.StdEncoding.EncodeToString(hash)
}

func RegisterNewUsers(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context){
		body, err := io.ReadAll(c.Request.Body)
        if err != nil {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
            return
        }

		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		var userSlice []models.User

		if err := c.BindJSON(&userSlice); err != nil {
			var singleUser models.User
			if err := c.BindJSON(&singleUser); err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}
			userSlice = append(userSlice, singleUser)
		}

		requiredFields := utils.GetRequiredFields(reflect.TypeOf(models.User{}))
		for _, usr := range userSlice {
			if userJSON, err := json.Marshal(usr); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user struct from request."})
			} else {
				if !utils.CheckRequiredFields(userJSON, requiredFields) {
					return
				}
			}
		}
		
		coll := client.Database("gore").Collection("user")
		var validUsers []interface{}

		for _, usr := range userSlice {
			var existingUser models.User
			if err := coll.FindOne(context.TODO(), bson.M{"username": usr.Username}).Decode(&existingUser); err != nil {
				if err == mongo.ErrNoDocuments {
					usr.Password = string(HashAndSalt(usr.Password))
					validUsers = append(validUsers, usr)
				} else {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "An expected error occured.", "detail": err.Error()})
					return
				}
			} else {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "User already exists."})
				return
			}
		}

		if len(validUsers) > 0 {
			res, err := coll.InsertMany(context.TODO(), validUsers)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An error occured when inserting users.", "detail": err.Error()})
			} 
			c.IndentedJSON(http.StatusOK, gin.H{"message": "User(s) successfully added.", "ids": res.InsertedIDs})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No valid users to add."})
		}
	}
}

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
			newUser.Password = string(HashAndSalt(newUser.Password))
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
	if string(HashAndSalt(password)) != result.Password {
		return "Password provided does not match.", models.User{}, fmt.Errorf("incorrect password")
	}

	return "Login successful.\n\n", result, nil
}