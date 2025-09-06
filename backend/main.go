package main

import (
	controller "backend/controllers"
	mongo "backend/database"
	router "backend/router"
	"fmt"
)

func main() {
	mongo.Connect()
	controller.InitCollection()
	r := router.SetupRouter()
	fmt.Println("Server Started on PORT 8080")
	r.Run(":8080")
}
