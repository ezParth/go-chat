package controllers

import (
	auth "backend/auth"
	client "backend/database"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var userCollection *mongo.Collection

func InitCollection() {
	userCollection = client.Client.Database("go-chat").Collection("User")
}

func GetUsers(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("Hitting getUsers")
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var signupData struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.BindJSON(&signupData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	hashedPassword, err := auth.HashPassword(signupData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := bson.M{
		"username": signupData.Username,
		"email":    signupData.Email,
		"password": hashedPassword,
	}
	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	token, err := auth.GenerateJWTToken(signupData.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, bson.M{
		"token":   token,
		"message": "Signup successful",
		"success": true,
	})
}

func Login(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var LoginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&LoginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var user bson.M
	err := userCollection.FindOne(ctx, bson.M{"email": LoginData.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	passwordFromDB, ok := user["password"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password not found or invalid in DB"})
		return
	}

	fmt.Println("Hashed Password -> ", passwordFromDB)
	isPasswordCorrect := auth.CompareHashedPassword(passwordFromDB, LoginData.Password)
	if !isPasswordCorrect {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	token, err := auth.GenerateJWTToken(user["name"].(string))
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "token": token, "message": "Successfully LoggedIn"})
}
