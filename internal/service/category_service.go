package service

import (
	"context"
	"errors"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/repository"
)

type ProductCategoryService struct {
	repo repository.ProductCategoryRepository // A reference to the ProductCategoryRepository interface in order to access its methods.
}

func NewCategoryService(cat_Repo repository.ProductCategoryRepository) *ProductCategoryService {
	return &ProductCategoryService{repo: cat_Repo}
}

// A function to check validation of Category
func categoryValidation(category *domain.ProductCategory) error {
	if category.Title == "" {
		return errors.New("name of category is required")
	}
	return nil
}

func (s *ProductCategoryService) CreateCategory(ctx context.Context, category *domain.ProductCategory) error {
	// Validation
	err := categoryValidation(category)
	if err != nil {
		return err
	}

	// Passing context to repository
	return s.repo.Create(ctx, category)
}

func (s *ProductCategoryService) ListCategories(ctx context.Context) ([]*domain.ProductCategory, error) {
	return s.repo.List(ctx)
}

func (s *ProductCategoryService) GetCategoryID(ctx context.Context, title string) (string, error) {
	category, err := s.repo.GetByTitle(ctx, title)
	if err != nil {
		return "", nil
	}
	return category.ID, nil // Only returning categoryID since rest of metadata of category collection is not needed by product handler, it only wants the categoryID
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
	return s.repo.Delete(ctx, id)
}
