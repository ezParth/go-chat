package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	GroupName string             `bson"groupName", json:"groupName"`
}
