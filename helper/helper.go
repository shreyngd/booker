package helper

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB : This is helper function to connect mongoDB
func ConnectDB() *mongo.Collection {
	fmt.Print(os.Getenv("MONGO_URI"))
	clientOpts := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	client,err := mongo.Connect(context.TODO(),clientOpts)

	if(err!= nil){
		log.Fatal(err)
	}

	collection := client.Database("booker").Collection("books")
	return collection;

}
