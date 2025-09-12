package main

import (
	helper "backend/config"
	controller "backend/controllers"
	mongo "backend/database"
	router "backend/router"
	"fmt"
)

func main() {
	mongo.Connect()
	controller.InitCollection()
	controller.InitGroupCollection()
	controller.InitConversationCollection()
	hub := helper.CreateHub()
	r := router.SetupRouter(hub)
	fmt.Println("Server Started on PORT 8080")
	r.Run(":8080")
}
