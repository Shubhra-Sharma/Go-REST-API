package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Product struct {
	ID          bson.ObjectID ` bson:"_id,omitempty"`
	Name        string        ` bson:"name"`
	Description string        ` bson:"description"`
	Category    string        ` bson:"category"`
	Price       int           ` bson:"price"`
	Brand       string        ` bson:"brand"`
	Quantity    int           ` bson:"quantity"`
}
