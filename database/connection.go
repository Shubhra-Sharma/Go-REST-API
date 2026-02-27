package database

import (
	 "context"
    "fmt"
    "log"
    "time"
    "os"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
    "go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// connect to mongoDB
func Connect(ctx context.Context) (*mongo.Database,error) {
    // getting environment variables
	mongoURI := os.Getenv("MONGO_URI")
	dbname := os.Getenv("DBNAME")

	// creating context for connection
	connectionCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
	//client option
	clientOption := options.Client().ApplyURI(mongoURI)

	//connecting to mongoDB
	client, err := mongo.Connect(connectionCtx,clientOption)

	if err!=nil{
	  return nil,err
	}

	fmt.Println("MongoDB connection successful!")

	// ping to ensure connection is established
    pingCtx, pingCancel := context.WithTimeout(ctx, 2*time.Second)
    defer pingCancel()

    if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
        return nil,err
    }

    fmt.Println("MongoDB connection successful!")
    
	// reference to database
	db := client.Database(dbName)
    return db, nil
}