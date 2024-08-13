package main

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
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://localhost:27017"

var currentUser models.User

func loginLogic(username string, password string, client *mongo.Client) (string, models.User, error) {
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

func loginCall(client *mongo.Client) gin.HandlerFunc {
	return func (c *gin.Context)  {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&credentials); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		message, u, error := loginLogic(credentials.Username, credentials.Password, client)
		if error != nil {
			c.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": error, "user": u, "message": message})
		} else {
			currentUser = u
			userResponse := utils.ConvertUserFormat([]models.User{u})
			c.String(http.StatusOK, message)
			c.IndentedJSON(http.StatusOK, gin.H{"userDetails": userResponse})
		}
	}
}

func welcomeMessage(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to my crap-app!")
}

func personalisedWelcome(c *gin.Context) {
	name := c.Param("name")

	switch {
	case !utils.NameValidCheck(name):
		c.String(http.StatusNotAcceptable, "Sorry, you have provided an invalid name (make sure there are no numbers or special characters and try again) ¯\\_(ツ)_/¯")
	case name == "Cameron":
		c.String(http.StatusOK, "Of course I know him - he's me!")
	default:
		message := fmt.Sprintf("Howdy %v, welcome to my crap-app!", name)
		c.String(http.StatusOK, message)
	}
}

func connectMongoDB(uri string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return nil, err
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
}

func getUserList(client *mongo.Client) gin.HandlerFunc {
	return func (c *gin.Context)  {
		if currentUser.Username != "" {
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
			}
			userRepsonses := utils.ConvertUserFormat(userList)
			c.IndentedJSON(http.StatusFound, userRepsonses)
		} else {
			c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You are not logged in as a user with access to this resource."})
		}
	}
}

func registerNewUser(client *mongo.Client) gin.HandlerFunc {
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

		if existingUser := coll.FindOne(context.TODO(), bson.M{"username": newUser.Username}); existingUser.Err() != nil {
			if err == mongo.ErrNoDocuments {
				res, err := coll.InsertOne(context.TODO(), newUser)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
					return
				}

				c.IndentedJSON(http.StatusOK, gin.H{"inserted_id": res.InsertedID})
			} else {
				c.IndentedJSON(400, gin.H{"error": "An unexpected error occured.", "details": err})
			}
		}

		c.IndentedJSON(400, gin.H{"message": "A user with these details already exists on our system."})
	}
}

func main() {
	client, err := connectMongoDB(uri)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	
	router := gin.Default()
	router.GET("/", welcomeMessage)
	router.GET("/:name", personalisedWelcome)
	router.POST("/login", loginCall(client))
	router.GET("/db/users/list", getUserList(client))
	router.POST("/db/users/register", registerNewUser(client))

	router.Run("localhost:8080")
}