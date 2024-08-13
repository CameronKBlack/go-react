package main

import (
	"context"
	"fmt"
	"net/http"
	"unicode"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://localhost:27017"

type user struct {
    Username      string    `json:"username" bson:"username"`
    FirstName     string    `json:"first_name" bson:"first_name"`
    LastName      string    `json:"last_name" bson:"last_name"`
    DateOfBirth   string	`json:"date_of_birth" bson:"date_of_birth"`
    EmailAddress  string    `json:"email_address" bson:"email_address"`
    Password      string    `json:"-" bson:"password"`
}

func (u user) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{
		"username": "%s",
		"first_name": "%s",
		"last_name": "%s",
		"date_of_birth": "%s",
		"email_address": "%s"
	}`, u.Username, u.FirstName, u.LastName, u.DateOfBirth, u.EmailAddress)), nil
}

func loginLogic(username string, password string, client *mongo.Client) (string, user, error) {
	var result user
	userColl := client.Database("gore").Collection("user")

	err := userColl.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "User not found.", user{}, err
		}
		return "Failed to fetch user.", result, err
	}
	if password != result.Password {
		return "Password provided does not match.", user{}, fmt.Errorf("incorrect password")
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

		message, user, error := loginLogic(credentials.Username, credentials.Password, client)
		if error != nil {
			c.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": error, "user": user, "message": message})
		} else {
			currentUser = user
			c.String(http.StatusOK, message)
			c.IndentedJSON(http.StatusOK, gin.H{"userDetails": user})
		}
	}
}

func nameValidCheck(name string) bool {
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			return false
		}
	}

	return true
}

func welcomeMessage(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to my crap-app!")
}

func personalisedWelcome(c *gin.Context) {
	name := c.Param("name")

	switch {
	case !nameValidCheck(name):
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

var currentUser user

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

			c.IndentedJSON(http.StatusFound, users)
		} else {
			c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You are not logged in as a user with access to this resource."})
		}
	}
}

func registerNewUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser user
		if err := c.BindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		coll := client.Database("gore").Collection("user")

		res, err := coll.InsertOne(context.TODO(), newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"inserted_id": res.InsertedID})
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