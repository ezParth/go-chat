package router

import (
	"github.com/gin-gonic/gin"

	auth "backend/auth"
	controller "backend/controllers"
	ws "backend/ws"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/", auth.AuthMiddleware(), controller.GetUsers)
		userRoutes.POST("/login", controller.Login)
		userRoutes.POST("/signup", controller.Signup)
	}

	r.GET("/ws", ws.WsHandler)

	return r
}
