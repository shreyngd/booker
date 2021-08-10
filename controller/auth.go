package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/shreyngd/booker/db"
	"github.com/shreyngd/booker/helper"
	"github.com/shreyngd/booker/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

// Login handler for uname and password
func (c *Controller) Login(ctx *gin.Context) {
	ctxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"data": "Invalid json provided",
		})
	}

	foundUser, err := db.FindUserByEmail(ctxt, user.Email)
	defer cancel()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data": "error occured while checking for the email",
		})
	}
	passwordIsValid, msg := db.VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if !passwordIsValid {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, foundUser.User_id)
	updatedUser, err := db.UpdateUserByID(token, refreshToken, foundUser.User_id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data": "error occured while logging in",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": updatedUser,
	})
}

//CreateUser is the api used to tget a single user
func (con *Controller) Signup(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	count, err := db.FindUserCountByMail(*user.Email)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
		return
	}
	password := db.HashPassword(*user.Password)
	user.Password = &password
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken
	resultInsertionNumber, insertErr := db.InsertUser(&user)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user item was not created"})
		return
	}
	c.JSON(http.StatusCreated, resultInsertionNumber)
}
