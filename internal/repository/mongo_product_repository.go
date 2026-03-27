package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/repository/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongoDB specific implementation of ProductRepository
type MongoProductRepository struct {
	collection *mongo.Collection
}

// Initializing MongoProductRepository to implement all the methods of ProductRepository.
func NewMongoProductRepository(client *mongo.Client, dbName string, collectionName string) *MongoProductRepository {
	db := client.Database(dbName)
	return &MongoProductRepository{collection: db.Collection(collectionName)}
}

// Inserting new product into product collection stored in database.
func (r *MongoProductRepository) Create(ctx context.Context, product *domain.Product) error {
	dbProduct, err := ToMongoProduct(product)
	if err != nil {
		return err
	}
	_, err = r.collection.InsertOne(ctx, dbProduct)
	return err
}

// Extracting the product with the particular id from the database
func (r *MongoProductRepository) Get(ctx context.Context, id string) (*domain.Product, error) {
	objectID, err := bson.ObjectIDFromHex(id) // converting string id to ObjectID which is what is recognized by MongoDB.
	if err != nil {
		return nil, err
	}

	var product models.Product
	filter := bson.M{"_id": objectID} // bson.M{} is a  map used to create MongoDB queries/filters, it is shorthand for "type M map[string]interface{}"
	err = r.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}
	return ToDomainProduct(&product), nil
}

// Get products by category
func (r *MongoProductRepository) GetByCategory(ctx context.Context, id string) ([]*domain.Product, error) {
	categoryID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID format: %w", err)
	}

	// This cursor finds all the products with the specific category_id
	cursor, err := r.collection.Find(ctx, bson.M{"category_id": categoryID})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products by category: %w", err)
	}
	defer cursor.Close(ctx)

	// Extracting products from cursor into dbProducts
	var dbProducts []*models.Product
	if err = cursor.All(ctx, &dbProducts); err != nil {
		return nil, err
	}

	// Converting slice of repo model products to domain model products
	products := make([]*domain.Product, len(dbProducts))
	for i, val := range dbProducts {
		products[i] = ToDomainProduct(val)
	}
	return products, nil
}

// Get all the products
func (r *MongoProductRepository) List(ctx context.Context) ([]*domain.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}) // Cursor is like a pointer that lets you iterate through multiple documents returned by a query.
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx) // we need to close the cursor after completion of function to prevent memory leak.

	var dbProducts []*models.Product // sending reference to slice in place of slice to save memory.
	if err = cursor.All(ctx, &dbProducts); err != nil {
		return nil, err
	}

	// Converting slice of repo models to slice of domain models
	products := make([]*domain.Product, len(dbProducts))
	for i, val := range dbProducts {
		products[i] = ToDomainProduct(val)
	}
	return products, nil
}

// Update a particular record in database with the help of its ID
func (r *MongoProductRepository) Update(ctx context.Context, id string, product *domain.Product) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid product ID format: %w", err)
	}

	// Here currently i am using the domain model directly to update fields in the mongo model since all the fields present in domain are also present in mongo model, the mongo model does not contain metadata like updated_at, so that's why in this specific scenario, directly using domain model works
	update := bson.M{
		"$set": bson.M{
			"name":        product.Name, // do not mutate or manipulate ID
			"description": product.Description,
			"category":    product.Category,
			"price":       product.Price,
			"brand":       product.Brand,
			"quantity":    product.Quantity,
		},
	}
	// this is an update instruction for mongoDB using $set operator.
	// $set updates all the fields with new values, the values of the rest of the fields remain unchanged.

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("product not found")
	}
	return nil
}

// Delete a particular record from db
func (r *MongoProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("product not found")
	}
	return nil
}

// Product == repository model, domain.Product == Domain Model
