package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct{
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Category string `json:"category" bson:"category"`
	Price int `json:"price" bson:"price"`
	Brand string `json:"brand" bson:"brand"`
	Quantity int `json:"quantity" bson:"quantity"`
}
