package router

import (
	"github.com/gin-gonic/gin"

	auth "backend/auth"
	controller "backend/controllers"
	helper "backend/helper"
	ws "backend/ws"
)

func SetupRouter(hub *helper.Hub) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(auth.CORSMiddleware())
	r.Use(func(ctx *gin.Context) {
		ctx.Set("hub", hub)
		ctx.Next()
	})

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/", auth.AuthMiddleware(), controller.GetUsers)
		userRoutes.POST("/login", controller.Login)
		userRoutes.POST("/signup", controller.Signup)
	}

	groupRoutes := r.Group("/group")
	{
		groupRoutes.POST("/create", auth.AuthMiddleware(), controller.CreateGroup)
		// groupRoutes.DELETE("/delete", auth.AuthMiddleware(), controller.CreateGroup())
		groupRoutes.POST("/join", auth.AuthMiddleware(), controller.JoinGroup)
		groupRoutes.GET("/getGroups", auth.AuthMiddleware(), controller.GetGroupsByUser)
	}

	r.GET("/ws", ws.WsHandler)

	return r
}
