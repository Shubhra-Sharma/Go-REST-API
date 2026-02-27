package repository
import(
	"os"
	"context"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/bson/primitive"
    "go.mongodb.org/mongo-driver/v2/mongo"
)
var collectionName string= os.Getenv("COLLECTION_NAME")

type MongoProductRepository struct {
	collection *mongo.Collection
}

// Initializing MongoProductRepository to implement all the methods of ProductRepository.
func NewMongoProductRepository(db *mongo.Database, collectionName string) *MongoProductRepository {
	return &MongoProductRepository{collection: db.Collection(collectionName)}
}

//Inserting new product into product collection stored in database.
func (r *MongoProductRepository) Create(ctx context.Context, product *domain.Product) error {
    result, err := r.collection.InsertOne(ctx, product)
    return err
}

//Extracting the product with the particular id from the database
func (r *MongoProductRepository) Get(ctx context.Context, id string) (*domain.Product, error) {
    objectID, err := primitive.ObjectIDFromHex(id) // converting string id to ObjectID which is what is recognized by MongoDB.
    if err != nil {
        return nil, err
    }
    
    var product domain.Product
	filter:= bson.M{"_id": objectID}  //bson.M{} is a  map used to create MongoDB queries/filters, it is shorthand for "type M map[string]interface{}"
    err = r.collection.FindOne(ctx, filter).Decode(&product)
    if err != nil {
        return nil, err
    }
    return &product, nil
}


// Get all the products
func (r *MongoProductRepository) List(ctx context.Context) ([]*domain.Product, error) {
    cursor, err := r.collection.Find(ctx, bson.M{}) //Cursor is like a pointer that lets you iterate through multiple documents returned by a query.
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx) // we need to close the cursor after completion of function to prevent memory leak.
    
    var products []*domain.Product // sending reference to slice in place of slice to save memory.
    if err = cursor.All(ctx, &products); err != nil {
        return nil, err
    }
    return products, nil
}

// Update a particular record in database with the help of its id
func (r *MongoProductRepository) Update(ctx context.Context, id string, product *domain.Product) error {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    
	update := bson.M{"$set": product}
// this is an update instruction for mongoDB using $set operator. 
// $set updates all the fields with new values, the values of the rest of the fields remain unchanged.

    _, err = r.collection.UpdateOne(
		 ctx,
		 bson.M{"_id": objectID}, 
		 update
		)
    return err
}

// Delete a particular record from db
func (r *MongoProductRepository) Delete(ctx context.Context, id string) error {
    objectID, err := primitive.ObjectIDFromHex(id) 
    if err != nil {
        return err
    }
    
    _, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
    return err
}