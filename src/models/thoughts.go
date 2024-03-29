package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Thought struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Text      string             `json:"text" validate:"required"`
	Anonymous bool               `json:"anonymous"`
	Replies   []Reply
}
