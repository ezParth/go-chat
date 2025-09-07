package main

import (
	controller "backend/controllers"
	mongo "backend/database"
	"backend/helper"
	router "backend/router"
	"fmt"
)

func main() {
	mongo.Connect()
	controller.InitCollection()
	controller.InitGroupCollection()
	hub := helper.CreateHub()
	r := router.SetupRouter(hub)
	fmt.Println("Server Started on PORT 8080")
	r.Run(":8080")
}
