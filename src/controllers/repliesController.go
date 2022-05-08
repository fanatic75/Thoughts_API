package controllers

import (
	"context"
	"fmt"
	"net/http"
	"thoughts-api/src/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddReplies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var reply models.Reply
		if err := c.BindJSON(&reply); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(reply)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "message": validationErr.Error()})
			return
		}

		reply.ID = primitive.NewObjectID()
		reply.Username = c.MustGet("username").(string)
		thoughtID, err := primitive.ObjectIDFromHex(c.Param("thoughtID"))
		if err != nil {
			msg := fmt.Sprintf("Thought ID is not correct")
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "message": msg})
			return
		}
		var documentToInsert = bson.M{"thoughts.$.replies": reply}
		err = models.UserCollection.FindOneAndUpdate(ctx,
			bson.M{"thoughts._id": thoughtID},
			bson.M{"$push": documentToInsert},
		).Decode(&bson.M{})
		if err != nil {
			msg := fmt.Sprintf("Reply was not created")
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": documentToInsert, "message": "Reply was added"})

	}
}

func DeleteReplies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		replyID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			msg := fmt.Sprintf("Reply id is not proper")
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": msg})
			return
		}
		thoughtID, err := primitive.ObjectIDFromHex(c.Param("thoughtID"))
		if err != nil {
			msg := fmt.Sprintf("Reply id is not proper")
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": msg})
			return
		}
		var userObject models.User
		err = models.UserCollection.FindOneAndUpdate(ctx,
			bson.M{"thoughts._id": thoughtID},
			bson.M{"$pull": bson.M{"thoughts.$.replies": bson.M{"_id": replyID, "username": c.MustGet("username").(string)}}},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
			options.FindOneAndUpdate().SetProjection(bson.M{"thoughts.replies.username": 0}),
		).Decode(&userObject)
		if err != nil {
			msg := fmt.Sprintf("No Such Thought")
			c.JSON(http.StatusNotFound,
				gin.H{"success": false, "data": nil, "message": msg})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": userObject.Thoughts, "message": "Reply was deleted"})

	}
}
