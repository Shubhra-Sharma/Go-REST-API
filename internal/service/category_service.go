package service

import (
	"context"
	"errors"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/repository"
)

type ProductCategoryService struct {
	repo        repository.ProductCategoryRepository // A reference to the ProductCategoryRepository interface in order to access its methods.
	productRepo repository.ProductRepository         //
}

func NewCategoryService(categoryRepo repository.ProductCategoryRepository, productRepository repository.ProductRepository) *ProductCategoryService {
	return &ProductCategoryService{repo: categoryRepo, productRepo: productRepository}
}

// A function to check validation of Category
func categoryValidation(category *domain.ProductCategory) error {
	if category.Title == "" {
		return errors.New("name of category is required")
	}
	return nil
}

func (s *ProductCategoryService) CreateCategory(ctx context.Context, category *domain.ProductCategory) (*domain.ProductCategory, error) {
	// Validation
	err := categoryValidation(category)
	if err != nil {
		return nil, err
	}

	// Passing context to repository
	return s.repo.Create(ctx, category)
}

func (s *ProductCategoryService) GetByID(ctx context.Context, id string) (*domain.ProductCategory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductCategoryService) ListCategories(ctx context.Context) ([]*domain.ProductCategory, error) {
	return s.repo.List(ctx)
}

func (s *ProductCategoryService) UpdateCategory(ctx context.Context, id string, category *domain.ProductCategory) error {
	// Validation
	err := categoryValidation(category)
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, id, category)
}

func (s *ProductCategoryService) DeleteCategory(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	// If category is successfully deleted from category collection, delete all the products with that particular category from product collection.
	return s.productRepo.DeleteByCategoryID(ctx, id)
}
