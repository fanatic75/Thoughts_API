package models

import (
	"thoughts-api/src/database"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Username   *string   `json:"username" validate:"required"`
	Password   *string   `json:"password" validate:"required,min=6"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

var UserCollection *mongo.Collection

func init() {
	UserCollection = database.OpenCollection(database.Client, "users")
}
