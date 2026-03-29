package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Shubhra-Sharma/Go-REST-API/database"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	// Loading env variables
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DBNAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	catcollectionName := os.Getenv("CATEGORY_COLLECTION_NAME")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is required")
	}
	if dbName == "" {
		log.Fatal("DBNAME environment variable is required")
	}
	if collectionName == "" {
		collectionName = "products"
	}
	if catcollectionName == "" {
		catcollectionName = "category"
	}

	client, err := database.Connect(ctx, mongoURI, dbName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := client.Database(dbName)
	productCollection := db.Collection(collectionName)
	catCollection := db.Collection(catcollectionName)

	// finding all the distinct category values
	var categoryNames []any
	result := productCollection.Distinct(ctx, "category", bson.M{
		"category": bson.M{"$exists": true},
	})

	// Decoding the result
	err = result.Decode(&categoryNames)
	if err != nil {
		log.Fatal("Error in decoding distinct categories:", err)
	}

	fmt.Printf("There are %d unique categories\n", len(categoryNames))

	// This categoryMap will store categoryNames and their respective categoryID
	categoryMap := make(map[string]bson.ObjectID)

	for _, category := range categoryNames {
		categoryName, ok := category.(string)
		if !ok || categoryName == "" {
			continue
		}

		// Check if any product category already exists in category collection
		var existingCategory bson.M
		err := catCollection.FindOne(ctx, bson.M{"title": categoryName}).Decode(&existingCategory)

		// If category does not exist in the category collection, make a categoryDoc and then append it to the category collection
		if err == mongo.ErrNoDocuments {
			categoryID := bson.NewObjectID()
			categoryDoc := bson.M{
				"_id":         categoryID,
				"title":       categoryName,
				"description": "Migrated category",
			}

			_, err := catCollection.InsertOne(ctx, categoryDoc)
			if err != nil {
				log.Printf("Could not create category '%s': %v", categoryName, err)
				continue
			}

			categoryMap[categoryName] = categoryID

		} else if err != nil {
			log.Printf("An error occurred while checking category :%v", err)
			continue
		} else {
			// if category already exists in collection, append its ID to the map
			categoryID := existingCategory["_id"].(bson.ObjectID)
			categoryMap[categoryName] = categoryID
		}
	}

	// Now append category_id for each category to product collection
	updatedProducts := 0
	for categoryName, categoryID := range categoryMap {
		// First find those products with the categoryName and no categoryID field
		filter := bson.M{
			"category":    categoryName,
			"category_id": bson.M{"$exists": false},
		}

		// Set category_id
		update := bson.M{
			"$set": bson.M{
				"category_id": categoryID,
			},
		}

		// Update and append category_id to product Collection
		result, err := productCollection.UpdateMany(ctx, filter, update)
		if err != nil {
			log.Printf("Could not append categoryID for '%s': %v", categoryName, err)
			continue
		}
		updatedProducts += int(result.ModifiedCount)
	}

	fmt.Printf("Products updated: %d\n", updatedProducts)
	fmt.Println("\nCategory_ID migration completed successfully")
}
