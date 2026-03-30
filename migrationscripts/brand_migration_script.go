package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Shubhra-Sharma/Go-REST-API/database"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	// Loading env variables
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DBNAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is required")
	}
	if dbName == "" {
		log.Fatal("DBNAME environment variable is required")
	}
	if collectionName == "" {
		collectionName = "products"
	}

	// Connecting to MongoDB
	client, err := database.Connect(ctx, mongoURI, dbName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := client.Database(dbName)
	productCollection := db.Collection(collectionName)

	// First this script finds products with empty brands
	filter := bson.M{
		"$or": []bson.M{
			{"brand": nil}, // brand is nil
			{"brand": ""},  // brand is empty string
		},
	}
	// brand field of all the products with no brand will be set to NA
	update := bson.M{
		"$set": bson.M{
			"brand": "NA",
		},
	}

	result, err := productCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Fatal("Failed to update products:", err)
	}
	fmt.Println("Brand migration successfully completed")
	if result.ModifiedCount == 0 {
		fmt.Println("There were 0 products with no brand")
	}
}
