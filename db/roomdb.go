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
		var room models.InterviewRoom
		err := cur.Decode(&room)
		if err != nil {
			log.Fatal(err)
		}
		rooms = append(rooms, room)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return rooms, nil
}
