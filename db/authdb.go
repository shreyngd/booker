package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shreyngd/booker/helper"
	"github.com/shreyngd/booker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

var collectionUser = helper.ConnectDB().Collection("users")

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

//VerifyPassword checks the input password while verifying it with the password in the DB.
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

//FindUserCountByMail to verify
func FindUserCountByMail(email string) (int64, error) {
	count, err := collectionUser.CountDocuments(context.TODO(), bson.M{"email": email})
	if err != nil {
		return count, err
	}
	return count, nil
}

//Insert user into DB
func InsertUser(u *models.User) (interface{}, error) {
	result, err := collectionUser.InsertOne(context.TODO(), u)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, err
}

//FindUser by email search
func FindUserByEmail(ctx context.Context, email *string) (models.User, error) {
	var foundUser models.User
	err := collectionUser.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
	fmt.Println(foundUser)
	if err != nil {
		return foundUser, err
	}
	return foundUser, err
}

// Update user
func UpdateUserByID(token string, refresh string, id string) (models.User, error) {
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{"token", token})
	updateObj = append(updateObj, bson.E{"refresh_token", refresh})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})
	upsert := true
	filter := bson.M{"user_id": id}
	opt := options.FindOneAndUpdateOptions{
		Upsert: &upsert,
	}
	opt.SetReturnDocument(options.After)
	var updatedUser models.User

	err := collectionUser.FindOneAndUpdate(
		context.TODO(),
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	).Decode(&updatedUser)

	fmt.Println(updatedUser)

	return updatedUser, err
}

func UpdateUserByIDAndGoogleToken(token string, refresh string, gToken *oauth2.Token, id string) (models.User, error) {
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{"token", token})
	updateObj = append(updateObj, bson.E{"refresh_token", refresh})
	updateObj = append(updateObj, bson.E{"googletoken", *gToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})
	upsert := true
	filter := bson.M{"user_id": id}
	opt := options.FindOneAndUpdateOptions{
		Upsert: &upsert,
	}
	opt.SetProjection(bson.D{
		{"email", 1},
		{"token", 1},
		{"googletoken", 1},
	})

	opt.SetReturnDocument(options.After)
	var updatedUser models.User

	err := collectionUser.FindOneAndUpdate(
		context.TODO(),
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	).Decode(&updatedUser)

	fmt.Println(updatedUser)

	return updatedUser, err
}
