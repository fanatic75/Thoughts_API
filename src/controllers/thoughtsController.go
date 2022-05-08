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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddThought() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var thought models.Thought
		if err := c.BindJSON(&thought); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		thought.ID = primitive.NewObjectID()
		thought.Replies = make([]models.Reply, 0)
		validationErr := validate.Struct(thought)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "message": validationErr.Error()})
			return
		}
		var documentToInsert = bson.M{"thoughts": thought}
		err := models.UserCollection.FindOneAndUpdate(ctx,
			bson.M{"username": c.MustGet("username").(string)},
			bson.M{"$push": documentToInsert},
		).Decode(&bson.M{})
		if err != nil {
			msg := fmt.Sprintf("Thought was not created")
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": documentToInsert, "message": "Thought was added"})

	}
}

func ListMyThoughts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		username := c.MustGet("username").(string)
		var matchPipeline = bson.D{
			{Key: "$match", Value: bson.M{
				"username": username,
			}},
		}

		var setAnonymousQuery = bson.D{{Key: "$set", Value: bson.M{
			"thoughts": bson.M{
				"$map": bson.M{
					"input": "$thoughts",
					"as":    "thought",
					"in": bson.M{
						"$mergeObjects": bson.A{
							"$$thought",
							bson.M{
								"replies": bson.M{
									"$map": bson.M{
										"input": "$$thought.replies",
										"as":    "reply",
										"in": bson.M{
											"$mergeObjects": bson.A{
												"$$reply",
												bson.M{
													"username": bson.M{
														"$cond": bson.M{
															"if":   "$$reply.anonymous",
															"then": "",
															"else": "$$reply.username",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}}}}

		cur, err := models.UserCollection.Aggregate(ctx, mongo.Pipeline{matchPipeline, setAnonymousQuery})
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": "Something wrong happened"})
			return
		}

		var userObjects []models.User
		if err = cur.All(ctx, &userObjects); err != nil {
			msg := fmt.Sprintf("User does not exist so thoughts are not there")
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": msg})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": userObjects[0].Thoughts, "message": "List of your thoughts"})

	}
}

func ListOtherUserThoughts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		username := c.Param("username")
		var matchPipeline = bson.D{
			{Key: "$match", Value: bson.M{
				"username": username,
			}},
		}

		var setAnonymousQuery = bson.D{{Key: "$set", Value: bson.M{
			"thoughts": bson.M{
				"$map": bson.M{
					"input": "$thoughts",
					"as":    "thought",
					"in": bson.M{
						"$mergeObjects": bson.A{
							"$$thought",
							bson.M{
								"replies": bson.M{
									"$map": bson.M{
										"input": "$$thought.replies",
										"as":    "reply",
										"in": bson.M{
											"$mergeObjects": bson.A{
												"$$reply",
												bson.M{
													"username": bson.M{
														"$cond": bson.M{
															"if":   "$$reply.anonymous",
															"then": "",
															"else": "$$reply.username",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}}}}

		var filterAnonymousPipeline = bson.D{{
			Key: "$project", Value: bson.M{
				"thoughts": bson.M{
					"$filter": bson.M{
						"input": "$thoughts",
						"as":    "thought",
						"cond": bson.M{
							"$eq": bson.A{
								"$$thought.anonymous",
								false,
							},
						},
					},
				},
			},
		}}

		//filtering out the thoughts which are anonymous
		cur, err := models.UserCollection.Aggregate(ctx, mongo.Pipeline{matchPipeline, filterAnonymousPipeline, setAnonymousQuery})
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": "Something wrong happened"})
			return
		}

		var userObjects []models.User
		if err = cur.All(ctx, &userObjects); err != nil || len(userObjects) == 0 {
			msg := fmt.Sprintf("No Thoughts for the specified user")
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": userObjects[0].Thoughts, "message": "List of thoughts"})

	}
}

func DeleteThoughts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		thoughtID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			msg := fmt.Sprintf("Thought id is not proper")
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": msg})
			return
		}
		var userObject models.User
		err = models.UserCollection.FindOneAndUpdate(ctx,
			bson.M{"username": c.MustGet("username").(string), "thoughts._id": thoughtID},
			bson.M{"$pull": bson.M{"thoughts": bson.M{"_id": thoughtID}}},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&userObject)
		if err != nil {
			msg := fmt.Sprintf("No Such Thought")
			c.JSON(http.StatusNotFound,
				gin.H{"success": false, "data": nil, "message": msg})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": userObject.Thoughts, "message": "Thought was deleted"})

	}
}
