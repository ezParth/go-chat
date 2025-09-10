package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	client "backend/database"
	"backend/models"
)

var groupCollection *mongo.Collection

func InitGroupCollection() {
	groupCollection = client.Client.Database("go-chat").Collection("Groups")
}

func CreateGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var groupData struct {
		GroupName string `json:"groupName"`
		Avatar    string `json:"avatar,omitempty"`
	}

	if err := c.BindJSON(&groupData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "success": false})
		return
	}

	group := models.Group{
		ID:        primitive.NewObjectID(),
		GroupName: groupData.GroupName,
		Admin: models.User{
			Username: username.(string),
		},
		Members: []models.User{
			{Username: username.(string)},
		},
		Avatar: groupData.Avatar,
	}

	_, err := groupCollection.InsertOne(ctx, group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group", "success": false})
		return
	}

	filter := bson.M{"username": username.(string)}
	update := bson.M{"$addToSet": bson.M{"groups": groupData.GroupName}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user with group", "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Group created successfully",
		"group":   group,
		"success": true,
	})
}

func JoinGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var groupData struct {
		GroupName string `json:"groupName"`
	}

	if err := c.BindJSON(&groupData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "success": false})
		return
	}

	var group models.Group
	err := groupCollection.FindOne(ctx, bson.M{"groupName": groupData.GroupName}).Decode(&group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group doesn't exist", "success": false})
		return
	}

	filter := bson.M{"username": username.(string)}
	update := bson.M{"$addToSet": bson.M{"groups": groupData.GroupName}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user with group", "success": false})
		return
	}

	groupUpdate := bson.M{"$addToSet": bson.M{"members": models.User{Username: username.(string)}}}
	_, err = groupCollection.UpdateOne(ctx, bson.M{"groupName": groupData.GroupName}, groupUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update group with user", "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Group joined successfully",
	})
}

func GetGroupsByUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username.(string)}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"groups":  user.Groups,
	})
}

func DeleteGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	username, exist := c.Get("username")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User doesn't exist", "success": false})
		return
	}

	var groupData struct {
		GroupName string `json:"groupName"`
	}

	if err := c.BindJSON(&groupData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err, "success": false})
	}

	var group models.Group
	if err := groupCollection.FindOne(ctx, bson.M{"groupname": groupData.GroupName}).Decode(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group does not exist", "success": false})
	}

	if group.Admin.Username != username.(string) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only admins can delete a group", "success": false})
	}

	_, err := groupCollection.DeleteOne(ctx, bson.M{"groupname": groupData.GroupName})
	if err != nil {
		fmt.Println("ERROR IN DELETING GROUP -> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete group", "success": false})
	}

	// will remove from user too later on -> might not be correct
	_, err = userCollection.UpdateMany(ctx, bson.M{"groups": groupData.GroupName}, bson.M{"$pull": bson.M{"groups": groupData.GroupName}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove group from users", "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted the group", "success": false})
}

func SaveGroupChat(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username, exist := c.Get("username")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized, please login", "success": false})
		return
	}

	var groupData struct {
		GroupName string `json:"groupName"`
		Message   string `json:"message"`
	}

	if err := c.BindJSON(&groupData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid data format", "success": false})
		return
	}

	chat := models.Chat{
		ID:       primitive.NewObjectID(),
		Sender:   username.(string),
		Receiver: groupData.GroupName, // receiver is group itself
		Message:  groupData.Message,
		Status:   "sent",
	}

	// push chat into group's messages array
	_, err := groupCollection.UpdateOne(
		ctx,
		bson.M{"groupname": groupData.GroupName},
		bson.M{"$push": bson.M{"messages": chat}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save chat", "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Chat saved successfully"})
}

func GetGroupChat(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	groupName := c.Param("groupName") // GET /groups/:groupName/chats

	var group models.Group
	err := groupCollection.FindOne(ctx, bson.M{"groupname": groupName}).Decode(&group)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"chats":   group.Messages,
	})
}

func GetGroupAvatar(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	groupName := c.Param("groupName")

	var group models.Group
	err := groupCollection.FindOne(ctx, bson.M{"groupname": groupName}).Decode(&group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group does not exist", "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully got the avatar",
		"avatar":  group.Avatar,
	})
}

func GetGroupMembersAndAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	groupName := c.Param("groupName")

	var group models.Group
	err := groupCollection.FindOne(ctx, bson.M{"groupname": groupName}).Decode(&group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group does not exist", "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"suceess": true, "message": "Successfully got members", "members": group.Members, "admin": group.Admin})
}
