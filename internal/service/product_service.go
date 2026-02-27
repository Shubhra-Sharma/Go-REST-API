package service

import (
    "context"
    "errors"
    "github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
    "github.com/Shubhra-Sharma/Go-REST-API/internal/repository"
)

type ProductService struct {
    repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
    return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *domain.Product) error {
    // Validation
    if product.Name=="" {
        return errors.New("product name is required")
    }
    if product.Price<=0 {
        return errors.New("product price must be greater than 0")
    }
    if product.Quantity<0 {
		return errors.New("Quantity must be greater than or equal to 0")
	}

	//Passing context to repository
    return s.repo.Create(ctx, product)
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
    return s.repo.Get(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
    return s.repo.List(ctx)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, product *domain.Product) error {
    // Validation
    if product.Name == "" {
        return errors.New("product name is required")
    }
    if product.Price <= 0 {
        return errors.New("product price must be greater than 0")
    }
    if product.Quantity<0 {
		return errors.New("Quantity must be greater than or equal to 0")
	}
    return s.repo.Update(ctx, id, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}
