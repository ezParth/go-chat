package router

import (
	"github.com/gin-gonic/gin"

	auth "backend/auth"
	"backend/controllers"
	controller "backend/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/", auth.AuthMiddleware(), controller.GetUsers)
		userRoutes.POST("/", auth.AuthMiddleware(), controller.CreateUser)
		userRoutes.GET("/:id", controllers.GetUserByID)
	}

	return r
}
