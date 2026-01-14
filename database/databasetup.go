package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client{
	var url string

	client, error := mongo.NewClient(options.Client().ApplyURI(url))

	if err!=nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	
	err = client.Ping(context.TODO(), nil)
	if err!=nil {
		log.Println("failed to connect to mongodb")
		return nil
	}

	fmt.Println("Connect mongodb successfuly!")
}

	var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Client {
	var collection *mongo.Collecions = client.Database("Ecommerce").Collection(collectionName)

	return collection
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Client {
	var collection *mongo.Collecions = client.Database("Ecommerce").Collection(collectionName)

	return collection
}