package db

import (
	"context"
	"log"

	"github.com/shreyngd/booker/helper"
	"github.com/shreyngd/booker/models"
	"go.mongodb.org/mongo-driver/bson"
)

var collectionRoom = helper.ConnectDB().Collection("rooms")

func CreateRoom(room *models.InterviewRoom) error {
	_, err := collectionRoom.InsertOne(context.TODO(), room)
	return err
}

func GetAllRooms() ([]models.InterviewRoom, error) {
	var rooms []models.InterviewRoom
	cur, err := collectionRoom.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var room models.InterviewRoom
		// & character returns the memory address of the following variable.
		err := cur.Decode(&room) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		rooms = append(rooms, room)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return rooms, nil
}
