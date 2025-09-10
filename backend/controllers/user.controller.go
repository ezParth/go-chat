package controllers

import (
	auth "backend/auth"
	client "backend/database"
	"backend/models"
	goType "backend/types"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func InitCollection() {
	userCollection = client.Client.Database("go-chat").Collection("User")
}

func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error fetching users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		log.Println("Error decoding users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode users"})
		return
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

	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: signupData.Username,
		Email:    signupData.Email,
		Password: hashedPassword,
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

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Signup successful",
		"success": true,
	})
}

// Login existing user
func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var user models.User
	// Allow login with either username or email
	filter := bson.M{
		"$or": []bson.M{
			{"username": loginData.Username},
			{"email": loginData.Username},
		},
	}

	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	isPasswordCorrect := auth.CompareHashedPassword(user.Password, loginData.Password)
	if !isPasswordCorrect {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	token, err := auth.GenerateJWTToken(user.Username)
	if err != nil {
		fmt.Println("Error in token generation -> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	fmt.Println("Logged in Successfully:", user.Username)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
		"message": "Successfully Logged In",
	})
}

func getUserByName(name string) *goType.GetUserByName {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	var user models.User

	if err := userCollection.FindOne(ctx, bson.M{"username": name}).Decode(&user); err != nil {
		return &goType.GetUserByName{Success: false, Message: "Cannot find User", User: nil}
	}

	return &goType.GetUserByName{Success: true, Message: "User Retrived Successfully", User: &user}
}

// package controllers

// import (
// 	auth "backend/auth"
// 	client "backend/database"
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type User struct {
// 	Username string `json:"username"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// var userCollection *mongo.Collection

// func InitCollection() {
// 	userCollection = client.Client.Database("go-chat").Collection("User")
// }

// func GetUsers(c *gin.Context) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	fmt.Println("Hitting getUsers")
// 	cursor, err := userCollection.Find(ctx, bson.M{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cursor.Close(ctx)
// 	var users []User
// 	if err = cursor.All(ctx, &users); err != nil {
// 		log.Fatal(err)
// 	}

// 	c.JSON(http.StatusOK, users)
// }

// func Signup(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	var signupData struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 		Email    string `json:"email"`
// 	}

// 	if err := c.BindJSON(&signupData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
// 		return
// 	}

// 	hashedPassword, err := auth.HashPassword(signupData.Password)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
// 		return
// 	}

// 	user := bson.M{
// 		"username": signupData.Username,
// 		"email":    signupData.Email,
// 		"password": hashedPassword,
// 	}
// 	_, err = userCollection.InsertOne(ctx, user)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
// 		return
// 	}

// 	token, err := auth.GenerateJWTToken(signupData.Username)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, bson.M{
// 		"token":   token,
// 		"message": "Signup successful",
// 		"success": true,
// 	})
// }

// func Login(c *gin.Context) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	var LoginData struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 	}

// 	if err := c.BindJSON(&LoginData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
// 		return
// 	}

// 	var user bson.M
// 	err := userCollection.FindOne(ctx, bson.M{"email": LoginData.Username}).Decode(&user)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
// 		return
// 	}

// 	passwordFromDB, ok := user["password"].(string)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "password not found or invalid in DB"})
// 		return
// 	}

// 	fmt.Println("USER ----> ", user)

// 	fmt.Println("Hashed Password -> ", passwordFromDB)
// 	isPasswordCorrect := auth.CompareHashedPassword(passwordFromDB, LoginData.Password)
// 	if !isPasswordCorrect {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
// 		return
// 	}

// 	token, err := auth.GenerateJWTToken(user["username"].(string))
// 	if err != nil {
// 		fmt.Println("Error in token generation -> ", err)
// 	}

// 	fmt.Println("Logged in Successfully")
// 	c.JSON(http.StatusOK, gin.H{"success": true, "token": token, "message": "Successfully LoggedIn"})
// }
