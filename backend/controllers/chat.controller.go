package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SaveUserChat(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	username, exist := c.Get("username")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Please login or create an account"})
		return
	}

	var chat struct {
		sender   string `json:"sender"`
		reciever string `json:"reciever"`
		message  string `json:"message"`
		time     string `json:"time"`
	}

	if err := c.BindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad format", "error": err, "success": false})
	}

}
