package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/shreyngd/booker/db"
	"github.com/shreyngd/booker/helper"
	"github.com/shreyngd/booker/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var validate = validator.New()

var (
	oauthConfGl = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
			"https://www.googleapis.com/auth/calendar"},
		Endpoint: google.Endpoint,
	}
)

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
	user.Role = "Interviewer"
	_, insertErr := db.InsertUser(&user)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user item was not created"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": user,
	})
}

//LoginGoogleURI := get googles uri for login
func (con *Controller) LoginGoogle(ctx *gin.Context) {
	state := helper.GetRandomState()
	u := oauthConfGl.AuthCodeURL(state)
	exp := time.Now().Add(24 * time.Hour)
	ctx.SetCookie("oauthstate", state, int(exp.Unix()), "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"redirectURI": u,
		},
	})
}

func (con *Controller) CallbackGoogle(ctx *gin.Context) {
	stateUrl := ctx.Query("state")
	stateCookie, err := ctx.Cookie("oauthstate")

	if err != nil {
		log.Println("invalid oauth google state")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data": "invalid state of request",
		})
		return
	}

	if stateUrl != stateCookie {
		log.Println("invalid oauth google state")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data": "invalid state of request",
		})
		return
	}

	token, err := oauthConfGl.Exchange(context.Background(), ctx.Query("code"))
	if err != nil {
		log.Fatalf("invalid code %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data": err,
		})
		return
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		log.Fatal("Get: " + err.Error() + "\n")
		ctx.JSON(http.StatusBadGateway, gin.H{
			"data": err,
		})
		return
	}
	defer resp.Body.Close()
	googleResponse := models.GoogleAuthResponse{}
	err = json.NewDecoder(resp.Body).Decode(&googleResponse)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"data": err,
		})
		return
	}

	u, err := db.FindUserByEmail(context.TODO(), &googleResponse.Email)
	if err != nil {
		var userCreate models.User
		userCreate.Email = &googleResponse.Email
		userCreate.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		userCreate.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		userCreate.ID = primitive.NewObjectID()
		userCreate.User_id = userCreate.ID.Hex()
		tokenSend, refreshToken, _ := helper.GenerateAllTokens(*userCreate.Email, userCreate.User_id)
		userCreate.Token = &tokenSend
		userCreate.Refresh_token = &refreshToken
		userCreate.GoogleToken = token
		userCreate.Role = "Interviewer"
		_, insertErr := db.InsertUser(&userCreate)
		if insertErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user item was not created"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"data": userCreate,
		})
		return
	}
	tokenSend, refreshToken, _ := helper.GenerateAllTokens(*u.Email, u.User_id)
	updatedUser, err := db.UpdateUserByIDAndGoogleToken(tokenSend, refreshToken, token, u.User_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data": "error occured while logging in",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": updatedUser,
	})
}
