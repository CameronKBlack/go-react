package routes

import (
	"go-react/backend/controllers"
	"go-react/backend/middlewares"
	"go-react/backend/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RouterSetup(client *mongo.Client) {
	router := gin.Default()

	router.GET("/", controllers.WelcomeMessage)
	router.GET("/:name", controllers.PersonalisedWelcome)
	router.POST("/login", controllers.LoginCall(client))
	router.POST("/db/users/register", services.RegisterNewUsers(client))

	authRoutes := router.Group("/")
	authRoutes.Use(middlewares.JWTMiddleware())
	{
		authRoutes.GET("/db/users/list", controllers.GetUserList(client))
	}
	router.Run("localhost:8080")
}