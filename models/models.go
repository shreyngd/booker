package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
)

//Create Struct
type Book struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Isbn   string             `json:"isbn,omitempty" bson:"isbn,omitempty"`
	Title  string             `json:"title" bson:"title,omitempty"`
	Author *Author            `json:"author" bson:"author,omitempty"`
}

type Author struct {
	FirstName string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Email         *string            `json:"email" validate:"email,required"`
	Password      *string            `json:"Password" validate:"required,min=6"`
	Token         *string            `json:"token"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
	Role          string             `json:"role"`
	GoogleToken   *oauth2.Token      `json:"googletoken"`
}

type SignedDetails struct {
	Email string
	Uid   string
	Role  string
	jwt.StandardClaims
}

type GoogleAuthResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type InterviewRoom struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       *string            `json:"name" validate:"name,required"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Active     bool               `json:"is_active"`
}

type AddRoomReply struct {
	RoomList  []string `json:"room_list"`
	AddedRoom string   `json:"added_room_name"`
}

type GlobalChannel struct {
	Channel chan string
}

var gb *GlobalChannel

func GetInstanceGlobal() *GlobalChannel {
	if gb == nil {
		gb = &GlobalChannel{
			Channel: make(chan string, 1000),
		}
	}
	return gb
}
