package controllers

import (
	client "backend/database"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var JWTKEY = []byte("My_JWT_Key")

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
var userCollection = client.Client.Database("go-chat").Collection("User")

func GetUsers(c *gin.Context) {
	defer cancel()
	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, users)
}

func Signup(c *gin.Context) {
	defer cancel()
	var user User
	user.Name = c.Param("username")
	user.Email = c.Param("email")
	user.Password = c.Param("password")
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, result)
}

func Login(c *gin.Context) {
	email := c.Param("email")
	var user bson.M
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusUnauthorized, bson.M{"error": err})
		return
	}

	c.JSON(http.StatusOK, user)
}
