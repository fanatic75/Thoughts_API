package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Reply struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Text      string             `json:"text" validate:"required"`
	Anonymous bool               `json:"anonymous"`
	Username  string             `bson:"username"`
}
