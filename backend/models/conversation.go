package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Users     []primitive.ObjectID `bson:"users" json:"users"` // two users stored here
	Messages  []Chat               `bson:"messages" json:"messages"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
}
