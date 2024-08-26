package controllers

import (
	"fmt"
	"go-react/backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WelcomeMessage(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to my crap-app!")
}

func PersonalisedWelcome(c *gin.Context) {
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
