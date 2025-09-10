package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatStatus string

const (
	StatusDelivered ChatStatus = "delivered"
	StatusReceived  ChatStatus = "received"
)

type Chat struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Message  string             `bson:"message" json:"message"`
	Sender   string             `bson:"sender" json:"sender"`
	Receiver string             `bson:"receiver" json:"receiver"`
	Status   ChatStatus         `bson:"status,omitempty" json:"status,omitempty"`
	Time     string             `bson:"time" json:"time"`
}

// Time     time.Time          `bson:"time" json:"time"` // will get date from frontend for time
// as it might take some time to save chats in backend
