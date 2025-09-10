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

func CreateGroupLogic(ctx context.Context, username, groupName, avatar string) (models.Group, error) {
	group := models.Group{
		ID:        primitive.NewObjectID(),
		GroupName: groupName,
		Admin:     models.User{Username: username},
		Members:   []models.User{{Username: username}},
		Avatar:    avatar,
	}

	_, err := groupCollection.InsertOne(ctx, group)
	if err != nil {
		return models.Group{}, err
	}

	filter := bson.M{"username": username}
	update := bson.M{"$addToSet": bson.M{"groups": groupName}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.Group{}, err
	}

	return group, nil
}

func JoinGroupLogic(ctx context.Context, username, groupName string) error {
	var group models.Group
	err := groupCollection.FindOne(ctx, bson.M{"groupName": groupName}).Decode(&group)
	if err != nil {
		return fmt.Errorf("group not found")
	}

	filter := bson.M{"username": username}
	update := bson.M{"$addToSet": bson.M{"groups": groupName}}
	if _, err := userCollection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	groupUpdate := bson.M{"$addToSet": bson.M{"members": models.User{Username: username}}}
	_, err = groupCollection.UpdateOne(ctx, bson.M{"groupName": groupName}, groupUpdate)
	return err
}

func GetGroupsByUserLogic(ctx context.Context, username string) ([]string, error) {
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user.Groups, nil
}

func DeleteGroupLogic(ctx context.Context, username, groupName string) error {
	var group models.Group
	if err := groupCollection.FindOne(ctx, bson.M{"groupname": groupName}).Decode(&group); err != nil {
		return fmt.Errorf("group not found")
	}

	if group.Admin.Username != username {
		return fmt.Errorf("only admins can delete group")
	}

	if _, err := groupCollection.DeleteOne(ctx, bson.M{"groupname": groupName}); err != nil {
		return err
	}

	_, err := userCollection.UpdateMany(ctx,
		bson.M{"groups": groupName},
		bson.M{"$pull": bson.M{"groups": groupName}},
	)
	return err
}

func SaveGroupChatLogic(ctx context.Context, username, groupName, message string) error {
	chat := models.Chat{
		ID:       primitive.NewObjectID(),
		Sender:   username,
		Receiver: groupName,
		Message:  message,
		Status:   "sent",
	}

	_, err := groupCollection.UpdateOne(
		ctx,
		bson.M{"groupname": groupName},
		bson.M{"$push": bson.M{"messages": chat}},
	)
	return err
}

func GetGroupChatLogic(ctx context.Context, groupName string) ([]models.Chat, error) {
	var group models.Group
	if err := groupCollection.FindOne(ctx, bson.M{"groupname": groupName}).Decode(&group); err != nil {
		return nil, err
	}
	return group.Messages, nil
}

func GetGroupAvatarLogic(ctx context.Context, groupName string) (string, error) {
	var group models.Group
	if err := groupCollection.FindOne(ctx, bson.M{"groupname": groupName}).Decode(&group); err != nil {
		return "", err
	}
	return group.Avatar, nil
}

func GetGroupMembersAndAdminLogic(ctx context.Context, groupName string) ([]models.User, models.User, error) {
	var group models.Group
	if err := groupCollection.FindOne(ctx, bson.M{"groupname": groupName}).Decode(&group); err != nil {
		return nil, models.User{}, err
	}
	return group.Members, group.Admin, nil
}

// ---------------- GIN HANDLERS ----------------

func CreateGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	username, _ := c.Get("username")
	var req struct {
		GroupName string `json:"groupName"`
		Avatar    string `json:"avatar,omitempty"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	group, err := CreateGroupLogic(ctx, username.(string), req.GroupName, req.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "group": group})
}

func JoinGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	username, _ := c.Get("username")
	var req struct {
		GroupName string `json:"groupName"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := JoinGroupLogic(ctx, username.(string), req.GroupName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "joined successfully"})
}

func GetGroupsByUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	username, _ := c.Get("username")
	groups, err := GetGroupsByUserLogic(ctx, username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch groups"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "groups": groups})
}

func DeleteGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	username, _ := c.Get("username")
	var req struct {
		GroupName string `json:"groupName"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err := DeleteGroupLogic(ctx, username.(string), req.GroupName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "deleted"})
}

func SaveGroupChat(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	username, _ := c.Get("username")
	var req struct {
		GroupName string `json:"groupName"`
		Message   string `json:"message"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err := SaveGroupChatLogic(ctx, username.(string), req.GroupName, req.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save chat"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetGroupChat(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	groupName := c.Param("groupName")
	chats, err := GetGroupChatLogic(ctx, groupName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "chats": chats})
}

func GetGroupAvatar(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	groupName := c.Param("groupName")
	avatar, err := GetGroupAvatarLogic(ctx, groupName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "avatar": avatar})
}

func GetGroupMembersAndAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	groupName := c.Param("groupName")
	members, admin, err := GetGroupMembersAndAdminLogic(ctx, groupName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "members": members, "admin": admin})
}

// package controllers

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"

// 	client "backend/database"
// 	"backend/models"
// )

// var groupCollection *mongo.Collection

// func InitGroupCollection() {
// 	groupCollection = client.Client.Database("go-chat").Collection("Groups")
// }

// func CreateGroup(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	username, exists := c.Get("username")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	var groupData struct {
// 		GroupName string `json:"groupName"`
// 		Avatar    string `json:"avatar,omitempty"`
// 	}

// 	if err := c.BindJSON(&groupData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "success": false})
// 		return
// 	}

// 	group := models.Group{
// 		ID:        primitive.NewObjectID(),
// 		GroupName: groupData.GroupName,
// 		Admin: models.User{
// 			Username: username.(string),
// 		},
// 		Members: []models.User{
// 			{Username: username.(string)},
// 		},
// 		Avatar: groupData.Avatar,
// 	}

// 	_, err := groupCollection.InsertOne(ctx, group)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group", "success": false})
// 		return
// 	}

// 	filter := bson.M{"username": username.(string)}
// 	update := bson.M{"$addToSet": bson.M{"groups": groupData.GroupName}}
// 	_, err = userCollection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user with group", "success": false})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Group created successfully",
// 		"group":   group,
// 		"success": true,
// 	})
// }

// func JoinGroup(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	username, exists := c.Get("username")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	var groupData struct {
// 		GroupName string `json:"groupName"`
// 	}

// 	if err := c.BindJSON(&groupData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "success": false})
// 		return
// 	}

// 	var group models.Group
// 	err := groupCollection.FindOne(ctx, bson.M{"groupName": groupData.GroupName}).Decode(&group)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Group doesn't exist", "success": false})
// 		return
// 	}

// 	filter := bson.M{"username": username.(string)}
// 	update := bson.M{"$addToSet": bson.M{"groups": groupData.GroupName}}
// 	_, err = userCollection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user with group", "success": false})
// 		return
// 	}

// 	groupUpdate := bson.M{"$addToSet": bson.M{"members": models.User{Username: username.(string)}}}
// 	_, err = groupCollection.UpdateOne(ctx, bson.M{"groupName": groupData.GroupName}, groupUpdate)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update group with user", "success": false})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "Group joined successfully",
// 	})
// }

// func GetGroupsByUser(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	username, exists := c.Get("username")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	var user models.User
// 	err := userCollection.FindOne(ctx, bson.M{"username": username.(string)}).Decode(&user)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"groups":  user.Groups,
// 	})
// }

// func DeleteGroup(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
// 	defer cancel()

// 	username, exist := c.Get("username")
// 	if !exist {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "User doesn't exist", "success": false})
// 		return
// 	}

// 	var groupData struct {
// 		GroupName string `json:"groupName"`
// 	}

// 	if err := c.BindJSON(&groupData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err, "success": false})
// 	}

// 	var group models.Group
// 	if err := groupCollection.FindOne(ctx, bson.M{"groupname": groupData.GroupName}).Decode(&group); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "group does not exist", "success": false})
// 	}

// 	if group.Admin.Username != username.(string) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "only admins can delete a group", "success": false})
// 	}

// 	_, err := groupCollection.DeleteOne(ctx, bson.M{"groupname": groupData.GroupName})
// 	if err != nil {
// 		fmt.Println("ERROR IN DELETING GROUP -> ", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete group", "success": false})
// 	}

// 	// will remove from user too later on -> might not be correct
// 	_, err = userCollection.UpdateMany(ctx, bson.M{"groups": groupData.GroupName}, bson.M{"$pull": bson.M{"groups": groupData.GroupName}})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove group from users", "success": false})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted the group", "success": false})
// }

// func SaveGroupChat(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	username, exist := c.Get("username")
// 	if !exist {
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized, please login", "success": false})
// 		return
// 	}

// 	var groupData struct {
// 		GroupName string `json:"groupName"`
// 		Message   string `json:"message"`
// 	}

// 	if err := c.BindJSON(&groupData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid data format", "success": false})
// 		return
// 	}

// 	chat := models.Chat{
// 		ID:       primitive.NewObjectID(),
// 		Sender:   username.(string),
// 		Receiver: groupData.GroupName,
// 		Message:  groupData.Message,
// 		Status:   "sent",
// 	}

// 	// push chat into group's messages array
// 	_, err := groupCollection.UpdateOne(
// 		ctx,
// 		bson.M{"groupname": groupData.GroupName},
// 		bson.M{"$push": bson.M{"messages": chat}},
// 	)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save chat", "success": false})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Chat saved successfully"})
// }

// func GetGroupChat(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	groupName := c.Param("groupName") // GET /groups/:groupName/chats

// 	var group models.Group
// 	err := groupCollection.FindOne(ctx, bson.M{"groupname": groupName}).Decode(&group)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Group not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"chats":   group.Messages,
// 	})
// }

// func GetGroupAvatar(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	groupName := c.Param("groupName")

// 	var group models.Group
// 	err := groupCollection.FindOne(ctx, bson.M{"groupname": groupName}).Decode(&group)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Group does not exist", "success": false})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "Successfully got the avatar",
// 		"avatar":  group.Avatar,
// 	})
// }

// func GetGroupMembersAndAdmin(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
// 	defer cancel()

// 	groupName := c.Param("groupName")

// 	var group models.Group
// 	err := groupCollection.FindOne(ctx, bson.M{"groupname": groupName}).Decode(&group)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Group does not exist", "success": false})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"suceess": true, "message": "Successfully got members", "members": group.Members, "admin": group.Admin})
// }
