package database

import (
	"context"
    "fmt"
    "time"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
    "go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// connect to mongoDB
func Connect(ctx context.Context, mongoURI string, dbName string) (*mongo.Database,error) {
	//client option
	clientOption := options.Client().ApplyURI(mongoURI)

	//connecting to mongoDB
	client, err := mongo.Connect(clientOption)

	if err!=nil{
	  return nil,err
	}
	// ping to ensure connection is established
    pingCtx, pingCancel := context.WithTimeout(ctx, 2*time.Second)
    defer pingCancel()

    if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		client.Disconnect(context.Background())  // Clean up if ping fails
        return nil,err
    }

    fmt.Println("MongoDB connection successful!")
    
	// reference to database
	db := client.Database(dbName)
    return db, nil
}