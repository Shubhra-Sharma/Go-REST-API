package service

import (
	"context"
	"errors"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/repository"
)

type ProductService struct {
	repo         repository.ProductRepository // A reference to the productRepository interface in order to access its methods.
	categoryRepo repository.ProductCategoryRepository
}

func NewProductService(product_Repo repository.ProductRepository, categoryRepository repository.ProductCategoryRepository) *ProductService {
	return &ProductService{repo: product_Repo, categoryRepo: categoryRepository}
}

// A function to check validation of product
func validation(product *domain.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	if product.Quantity < 0 {
		return errors.New("quantity must be greater than or equal to 0")
	}
	if product.Brand == "" {
		return errors.New("brand name is compulsory for all products")
	}
	return nil
}

func (s *ProductService) CreateProduct(ctx context.Context, product *domain.Product) error {
	// Validation
	err := validation(product)
	if err != nil {
		return err
	}

	// Keeping it seperate here since the Update function does not require validation of a category Name
	if product.Category == "" {
		return errors.New("category name is required")
	}

	// Calling categoryRepo to get category ID for product creation
	category, err := s.categoryRepo.GetByTitle(ctx, product.Category)
	if err != nil {
		return err
	}
	product.CategoryID = category.ID
	return s.repo.Create(ctx, product)
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.Get(ctx, id)
}

func (s *ProductService) GetByCategory(ctx context.Context, categoryTitle string) ([]*domain.Product, error) {
	// first fetching category ID from category Repository
	category, err := s.categoryRepo.GetByTitle(ctx, categoryTitle)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByCategory(ctx, category.ID)
}

func (s *ProductService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	return s.repo.List(ctx)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, product *domain.Product) error {
	// Validation
	err := validation(product)
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, id, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
