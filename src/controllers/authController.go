package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	helper "thoughts-api/src/helpers"
	"thoughts-api/src/models"

	"go.mongodb.org/mongo-driver/bson"
)

var validate = validator.New()

//HashPassword is used to encrypt the password
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

//VerifyPassword checks
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

//sign up user
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "message": validationErr.Error()})
			return
		}

		count, err := models.UserCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": "error occured while checking for the username"})

			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": "username ready exists"})

			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		user.Thoughts = make([]models.Thought, 0)

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		_, insertErr := models.UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": msg})

			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": nil, "message": "user signup sucess"})

	}
}

//Login user
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		err := models.UserCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Username == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		token, _ := helper.GenerateAllTokens(*foundUser.Username)

		err = models.UserCollection.FindOne(ctx, bson.M{"username": foundUser.Username}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{
			"token": token},
			"message": "return successfully"})

	}
}
