package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
)

func TestCreateCategory(t *testing.T) {
	testcases := []struct {
		name            string
		category        *domain.ProductCategory
		categoryRepoErr error
		want            string
	}{
		{ // default
			name:     "success",
			category: mockCategory(),
			want:     "",
		},
		{ // Testing validation logic
			name:     "missing title",
			category: &domain.ProductCategory{Description: "new category"},
			want:     "name of category is required",
		},
		{ // Checking if categoryService correctly returns repository errors or not
			name:            "Repository error",
			category:        mockCategory(),
			categoryRepoErr: errors.New("could not create new category."),
			want:            "could not create new category.",
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			categoryRepo := &mockCategoryRepo{
				create: func(ctx context.Context, category *domain.ProductCategory) (*domain.ProductCategory, error) {
					if tt.categoryRepoErr != nil {
						return nil, tt.categoryRepoErr
					}
					return category, nil
				},
			}
			categoryService := NewCategoryService(categoryRepo, &mockProductRepo{})
			result, err := categoryService.CreateCategory(context.Background(), tt.category)

			if tt.want == "" {
				if err != nil {
					t.Errorf("Expected no error, got %v.", err)
				}
				if result == nil {
					t.Errorf("Expected resultant category, got nil.")
				}
			} else {
				if err == nil || err.Error() != tt.want {
					t.Errorf("Expected error: %s, got %v", tt.want, err)
				}
			}
		})
	}
}

func TestUpdateCategory(t *testing.T) {
	testcases := []struct {
		name            string
		category        *domain.ProductCategory
		categoryRepoErr error
		want            string
	}{
		{ // default
			name:     "success",
			category: mockCategory(),
			want:     "",
		},
		{ // Testing validation logic
			name:     "missing title",
			category: &domain.ProductCategory{},
			want:     "name of category is required",
		},
		{ // Checking if categoryService correctly returns repository errors or not
			name:            "Repository error",
			category:        mockCategory(),
			categoryRepoErr: errors.New("could not update category."),
			want:            "could not update category.",
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			categoryRepo := &mockCategoryRepo{
				update: func(ctx context.Context, id string, category *domain.ProductCategory) error {
					return tt.categoryRepoErr
				},
			}
			categoryService := NewCategoryService(categoryRepo, &mockProductRepo{})
			err := categoryService.UpdateCategory(context.Background(), "123", tt.category)

			if tt.want == "" && err != nil {
				t.Errorf("Expected no error, got %v.", err)
			} else if tt.want != "" {
				if err == nil || err.Error() != tt.want {
					t.Errorf("Expected error: %s, got %v", tt.want, err)
				}
			}
		})
	}
}

// Delete fn of categoryService involves interaction with two repositories, thats why inluding it's test, is it right to do or not?
func TestDeleteCategory(t *testing.T) {
	tests := []struct {
		name            string
		categoryRepoErr error
		productRepoErr  error
		want            string
	}{
		{ // success case
			name: "success",
			want: "",
		},
		{ // category to be deleted was not found
			name:            "category not found",
			categoryRepoErr: errors.New("category not found"),
			want:            "category not found",
		},
		{
			// Category was deleted successfully but the operation for deleting the orphaned products failed
			name:           "Delete operation for orphaned products failed",
			productRepoErr: errors.New("failed to delete products"),
			want:           "failed to delete products",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryRepo := &mockCategoryRepo{
				delete: func(ctx context.Context, id string) error {
					return tt.categoryRepoErr
				},
			}

			productRepo := &mockProductRepo{
				deleteByCategory: func(ctx context.Context, id string) error {
					return tt.productRepoErr
				},
			}
			categoryService := NewCategoryService(categoryRepo, productRepo)
			err := categoryService.DeleteCategory(context.Background(), "123")

			if tt.want == "" && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
			if tt.want != "" {
				if err == nil || err.Error() != tt.want {
					t.Errorf("expected error '%s', got '%v'", tt.want, err)
				}
			}
		})
	}
}
