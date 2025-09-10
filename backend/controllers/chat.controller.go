package controllers

import (
	client "backend/database"
	"backend/models"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var conversationCollection *mongo.Collection

func InitConversationCollection() {
	conversationCollection = client.Client.Database("go-chat").Collection("Conversation")
}

func SaveUserChatToDB(sender string, receiver string, message string, t string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user1 := getUserByName(sender)
	user2 := getUserByName(receiver)

	if user1.User == nil || user2.User == nil {
		return errors.New("one or both users not found")
	}

	opts := options.Update().SetUpsert(true)

	newChat := models.Chat{
		ID:       primitive.NewObjectID(),
		Message:  message,
		Sender:   sender,
		Receiver: receiver,
		Status:   models.StatusDelivered,
		Time:     t,
	}

	filter := bson.M{
		"users": bson.M{"$all": []primitive.ObjectID{user1.User.ID, user2.User.ID}},
	}

	update := bson.M{
		"$push": bson.M{"messages": newChat},
		"$setOnInsert": bson.M{
			"_id":       primitive.NewObjectID(),
			"users":     []primitive.ObjectID{user1.User.ID, user2.User.ID},
			"createdAt": time.Now(),
		},
	}

	_, err := conversationCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func SaveUserChat(c *gin.Context) {
	username, exist := c.Get("username")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Please login or create an account"})
		return
	}

	var chat struct {
		Receiver string `json:"receiver"`
		Message  string `json:"message"`
		Time     string `json:"time"`
	}

	if err := c.BindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad format", "error": err.Error(), "success": false})
		return
	}

	err := SaveUserChatToDB(username.(string), chat.Receiver, chat.Message, chat.Time)
	if err != nil {
		if err.Error() == "one or both users not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save chat", "error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Chat saved successfully"})
}

// package controllers

// import (
// 	client "backend/database"
// 	"backend/models"
// 	"context"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// var conversationCollection *mongo.Collection

// func InitConversationCollection() {
// 	groupCollection = client.Client.Database("go-chat").Collection("Conversation")
// }

// func SaveUserChat(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
// 	defer cancel()

// 	username, exist := c.Get("username")
// 	if !exist {
// 		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Please login or create an account"})
// 		return
// 	}

// 	var chat struct {
// 		Receiver string `json:"receiver"`
// 		Message  string `json:"message"`
// 		Time     string `json:"time"`
// 	}

// 	if err := c.BindJSON(&chat); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad format", "error": err, "success": false})
// 		return
// 	}

// 	user1 := getUserByName(username.(string))
// 	user2 := getUserByName(chat.Receiver)

// 	if user1.User == nil || user2.User == nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"message": "One or both users not found",
// 		})
// 		return
// 	}

// 	opts := options.Update().SetUpsert(true)

// 	newChat := models.Chat{
// 		ID:       primitive.NewObjectID(),
// 		Message:  chat.Message,
// 		Sender:   username.(string),
// 		Receiver: chat.Receiver,
// 		Status:   models.StatusDelivered,
// 		Time:     chat.Time,
// 	}

// 	filter := bson.M{
// 		"users": bson.M{"$all": []primitive.ObjectID{user1.User.ID, user2.User.ID}},
// 	}

// 	update := bson.M{
// 		"$push": bson.M{"messages": newChat},
// 		"$setOnInsert": bson.M{
// 			"_id":       primitive.NewObjectID(),
// 			"users":     []primitive.ObjectID{user1.User.ID, user2.User.ID},
// 			"createdAt": time.Now(),
// 		},
// 	}

// 	_, err := conversationCollection.UpdateOne(ctx, filter, update, opts)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"message": "Failed to save chat",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "Chat saved successfully",
// 	})
// }
