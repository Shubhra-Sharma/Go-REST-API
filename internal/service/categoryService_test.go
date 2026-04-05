package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory(t *testing.T) {
	testcases := []struct {
		name            string
		category        *domain.ProductCategory
		categoryRepoErr error
		wantErrMsg      string
	}{
		{
			name:     "success",
			category: mockCategory(),
		},
		{
			name:       "missing title",
			category:   &domain.ProductCategory{Description: "new category"},
			wantErrMsg: "name of category is required",
		},
		{
			name:            "repository error",
			category:        mockCategory(),
			categoryRepoErr: errors.New("could not create new category."),
			wantErrMsg:      "could not create new category.",
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

			svc := NewCategoryService(categoryRepo, &mockProductRepo{})
			result, err := svc.CreateCategory(context.Background(), tt.category)

			if tt.wantErrMsg == "" {
				require.NoError(t, err)
				assert.NotNil(t, result)
			} else {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
				assert.Nil(t, result)
			}
		})
	}
}

func TestUpdateCategory(t *testing.T) {
	testcases := []struct {
		name            string
		category        *domain.ProductCategory
		categoryRepoErr error
		wantErrMsg      string
	}{
		{
			name:     "success",
			category: mockCategory(),
		},
		{
			name:       "missing title",
			category:   &domain.ProductCategory{},
			wantErrMsg: "name of category is required",
		},
		{
			name:            "repository error",
			category:        mockCategory(),
			categoryRepoErr: errors.New("could not update category."),
			wantErrMsg:      "could not update category.",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			categoryRepo := &mockCategoryRepo{
				update: func(ctx context.Context, id string, category *domain.ProductCategory) error {
					return tt.categoryRepoErr
				},
			}

			svc := NewCategoryService(categoryRepo, &mockProductRepo{})
			err := svc.UpdateCategory(context.Background(), "123", tt.category)

			if tt.wantErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
			}
		})
	}
}

// TestDeleteCategory covers two-repo interaction: category deletion + orphaned product cleanup
func TestDeleteCategory(t *testing.T) {
	tests := []struct {
		name            string
		categoryRepoErr error
		productRepoErr  error
		wantErrMsg      string
	}{
		{
			name: "success",
		},
		{
			name:            "category not found",
			categoryRepoErr: errors.New("category not found"),
			wantErrMsg:      "category not found",
		},
		{
			name:           "delete orphaned products failed",
			productRepoErr: errors.New("failed to delete products"),
			wantErrMsg:     "failed to delete products",
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

			svc := NewCategoryService(categoryRepo, productRepo)
			err := svc.DeleteCategory(context.Background(), "123")

			if tt.wantErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
			}
		})
	}
}

func TestListCategories(t *testing.T) {
	tests := []struct {
		name            string
		categoryRepoErr error
		mockResult      []*domain.ProductCategory
		wantErrMsg      string
	}{
		{
			name:       "success",
			mockResult: []*domain.ProductCategory{mockCategory(), mockCategory()},
		},
		{
			name:            "repository error",
			mockResult:      nil,
			categoryRepoErr: errors.New("failed to fetch categories"),
			wantErrMsg:      "failed to fetch categories",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryRepo := &mockCategoryRepo{
				list: func(ctx context.Context) ([]*domain.ProductCategory, error) {
					return tt.mockResult, tt.categoryRepoErr
				},
			}

			svc := NewCategoryService(categoryRepo, &mockProductRepo{})
			result, err := svc.ListCategories(context.Background())

			if tt.wantErrMsg == "" {
				require.NoError(t, err)
				assert.Equal(t, tt.mockResult, result) // verifies the service returns exactly what the repo returned
			} else {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
				assert.Nil(t, result)
			}
		})
	}
}
