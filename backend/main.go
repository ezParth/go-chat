package main

import (
	mongo "backend/database"
	router "backend/router"
)

func main() {
	mongo.Connect()
	r := router.SetupRouter()
	r.Run(":8080")
}
