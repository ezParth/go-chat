package router

import (
	"github.com/gin-gonic/gin"

	auth "backend/auth"
	controller "backend/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/", auth.AuthMiddleware(), controller.GetUsers)
		userRoutes.POST("/login", controller.Login)
		userRoutes.POST("/signup", controller.Signup)
	}

	return r
}
