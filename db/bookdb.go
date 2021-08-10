package db

import (
	"context"
	"encoding/json"
	"log"

	"github.com/shreyngd/booker/helper"
	"github.com/shreyngd/booker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// collection of books
var collection = helper.ConnectDB()


//List all books in db
func GetBooks() ([]models.Book,error){

	var books []models.Book

	cur, err := collection.Find(context.TODO(), bson.M{})
	
	if err != nil {
		return nil,err
	}

	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var book models.Book
		// & character returns the memory address of the following variable.
		err := cur.Decode(&book) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		books = append(books, book)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return books,nil
}

func PutBooks(books []models.Book) (int,error){
	for i := range books {
		books[i].ID = primitive.NewObjectID()
	}
	c,_ := json.Marshal(books)
	var b = []interface{}{}
	json.Unmarshal(c, &b)
	result, err := collection.InsertMany(context.TODO(),b)
	if err != nil {
		return -1, err
	}
	return len(result.InsertedIDs),nil
}