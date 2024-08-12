package main

import (
	"fmt"
	"net/http"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

type user struct {
	Username string `json:"username"`
	First_name string `json:"first_name"`
	Last_name string `json:"last_name"`
	Date_of_birth time.Time `json:"D.O.B"`
	Email_address string `json:"email_address"`
	Password string `json:"-"`
}

func (u user) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{
		"username": "%s",
		"first_name": "%s",
		"last_name": "%s",
		"D.O.B": "%s",
		"email_address": "%s"
	}`, u.Username, u.First_name, u.Last_name, u.Date_of_birth.Format("02-Jan-2006"), u.Email_address)), nil
}

var logins = map[string]user{
	"cblack51": {
		Username: "cblack51",
		First_name: "Cameron",
		Last_name: "Black",
		Date_of_birth: time.Date(1994, time.February, 8, 0, 0, 0, 0, time.UTC),
		Email_address: "cameronblack51@hotmail.co.uk",
		Password: "hunter2",
	},
}

func loginLogic(username string, password string) (string, user, error) {
	user, ok := logins[username]
	if !ok {
			return "", user, fmt.Errorf("user not found")
	}

	if user.Password != password {
			return "Incorrect password.", user, fmt.Errorf("incorrect password provided")
	}

	message := fmt.Sprintf("Login successful.\nWelcome, %v.\n", username)
	return message, user, nil
}

func loginCall(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	message, user, error := loginLogic(credentials.Username, credentials.Password)
	if error != nil {
		c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": error})
	} else {
		c.String(http.StatusOK, message)
		c.IndentedJSON(http.StatusOK, gin.H{"userDetails": user})
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

func main() {
	router := gin.Default()

	router.GET("/", welcomeMessage)
	router.GET("/:name", personalisedWelcome)
	router.POST("/login", loginCall)
	
	router.Run("localhost:8080")
}