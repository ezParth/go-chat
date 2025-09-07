package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Group struct {
	GroupName string
	Admin     User
	Members   []User
}

func CreateGroup(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	username, exists := c.Get("username")
	if !exists {
		fmt.Println("Username doesn't exist")
		c.JSON(http.StatusInternalServerError, bson.M{"ERROR": "Username Does Not Exist"})
		return
	}

}
