package models

import "go.mongodb.org/mongo-driver/v2/bson"

type ProductCategory struct {
	ID          bson.ObjectID ` bson:"_id,omitempty"`
	Title       string        ` bson:"title"`
	Description string        ` bson:"description"`
}
