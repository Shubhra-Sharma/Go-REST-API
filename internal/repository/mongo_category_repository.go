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

// MongoDB specific implementation of ProductCategoryRepository
type MongoProductCategoryRepository struct {
	collection *mongo.Collection
}

// Initializing MongoProductCategoryRepository to implement all the methods of ProductCategoryRepository.
func NewMongoProductCategoryRepository(client *mongo.Client, dbName string, categoryCollectionName string) *MongoProductCategoryRepository {
	db := client.Database(dbName)
	return &MongoProductCategoryRepository{collection: db.Collection(categoryCollectionName)}
}

// Inserting new category into category collection.
func (r *MongoProductCategoryRepository) Create(ctx context.Context, category *domain.ProductCategory) (*domain.ProductCategory, error) {
	dbCategory, err := ToMongoCategory(category)
	if err != nil {
		return nil, err
	}

	// Checking if the category already exists in the database to prevent duplication of same category
	var isDuplicate bson.M
	err = r.collection.FindOne(ctx, bson.M{"title": dbCategory.Title}).Decode(&isDuplicate)
	// If category already exists, do not create duplicate category
	if err == nil {
		return nil, fmt.Errorf("category '%s' already exists", dbCategory.Title)
	}
	if err != mongo.ErrNoDocuments {
		// A real error occurred
		return nil, err
	}

	// Inserting new category document to database
	_, err = r.collection.InsertOne(ctx, dbCategory)
	return ToDomainCategory(dbCategory), err
}

// Get all the categories in db
func (r *MongoProductCategoryRepository) List(ctx context.Context) ([]*domain.ProductCategory, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dbCategories []*models.ProductCategory // sending reference to slice in place of slice to save memory.
	if err = cursor.All(ctx, &dbCategories); err != nil {
		return nil, err
	}

	// Converting slice of repo models to slice of domain models
	categories := make([]*domain.ProductCategory, len(dbCategories))
	for i, val := range dbCategories {
		categories[i] = ToDomainCategory(val)
	}
	return categories, nil
}

func (r *MongoProductCategoryRepository) GetByID(ctx context.Context, id string) (*domain.ProductCategory, error) {
	objectID, err := bson.ObjectIDFromHex(id) // converting string id to ObjectID which is what is recognized by MongoDB.
	if err != nil {
		return nil, err
	}
	var category models.ProductCategory
	filter := bson.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		return nil, err
	}
	return ToDomainCategory(&category), nil
}

// This returns the category_id for a particular category name
func (r *MongoProductCategoryRepository) GetByTitle(ctx context.Context, title string) (*domain.ProductCategory, error) {
	var category models.ProductCategory
	filter := bson.M{"title": title}
	err := r.collection.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		return nil, err
	}
	return ToDomainCategory(&category), nil
}

// Update a particular record in collection with the help of its ID
func (r *MongoProductCategoryRepository) Update(ctx context.Context, id string, category *domain.ProductCategory) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid category ID format: %w", err)
	}
	update := bson.M{
		"$set": bson.M{
			"title":       category.Title, // do not mutate or manipulate ID
			"description": category.Description,
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
		return errors.New("category not found")
	}
	return nil
}

// Delete a particular record from db
func (r *MongoProductCategoryRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("category not found")
	}
	return nil
}

// models.Category == repository model, domain.Category == Domain Model
