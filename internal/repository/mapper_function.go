package repository

import (
	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/repository/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Mapping Domain's Product to Repository's Product
func ToMongoProduct(domainProduct *domain.Product) (*models.Product, error) {
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

	return &models.Product{
		ID:          objectID,
		Name:        domainProduct.Name,
		Description: domainProduct.Description,
		Category:    domainProduct.Category,
		Price:       domainProduct.Price,
		Brand:       domainProduct.Brand,
		Quantity:    domainProduct.Quantity,
	}, nil
}

// Mapping Repository's Product to Domain's Product
func ToDomainProduct(mongoProduct *models.Product) *domain.Product {
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

// Mapping Repository's ProductCategory to Domain's ProductCategory
func ToMongoCategory(domainCat *domain.ProductCategory) (*models.ProductCategory, error) {
	var objectID bson.ObjectID
	var err error

	// Handling objectID for cases where ID is empty and where it is not
	if domainCat.ID != "" {
		objectID, err = bson.ObjectIDFromHex(domainCat.ID)
		if err != nil {
			return nil, err
		}
	} else {
		objectID = bson.NewObjectID()
	}

	return &models.ProductCategory{
		ID:          objectID,
		Title:       domainCat.Title,
		Description: domainCat.Description,
	}, nil
}

// Mapping Repository's ProductCategory to Domain's ProdutCategory
func ToDomainCategory(mongoCat *models.ProductCategory) *domain.ProductCategory {
	return &domain.ProductCategory{
		ID:          mongoCat.ID.Hex(),
		Title:       mongoCat.Title,
		Description: mongoCat.Description,
	}
}
