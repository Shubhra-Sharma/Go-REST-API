package repository

import (
	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Mapping Domain Model to Repository Model
func ToMongoProduct(domainProduct *domain.Product) (*Product, error) {
	var objectID bson.ObjectID
	var err error

	// Handling objectID for cases where ID is empty and where it is not
	if domainProduct.ID != "" {
		objectID, err = bson.ObjectIDFromHex(domainProduct.ID)
		if err != nil {
			return nil, err
		}
	} else {
		objectID = bson.NewObjectID()
	}

	return &Product{
		ID:          objectID,
		Name:        domainProduct.Name,
		Description: domainProduct.Description,
		Category:    domainProduct.Category,
		Price:       domainProduct.Price,
		Brand:       domainProduct.Brand,
		Quantity:    domainProduct.Quantity,
	}, nil
}

// Mapping Repossitory Model to Domain Model
func ToDomainProduct(mongoProduct *Product) *domain.Product {
	return &domain.Product{
		ID:          mongoProduct.ID.Hex(),
		Name:        mongoProduct.Name,
		Description: mongoProduct.Description,
		Category:    mongoProduct.Category,
		Price:       mongoProduct.Price,
		Brand:       mongoProduct.Brand,
		Quantity:    mongoProduct.Quantity,
	}
}
