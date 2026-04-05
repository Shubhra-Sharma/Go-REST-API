package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	testcases := []struct {
		name            string
		product         *domain.Product
		categoryRepoErr error
		productRepoErr  error
		wantErrMsg      string // empty = expect success
	}{
		{
			name:    "success",
			product: mockProduct(),
		},
		{
			name:       "missing name",
			product:    &domain.Product{Price: 100, Quantity: 1, Brand: "Vivo", Category: "Device"},
			wantErrMsg: "product name is required",
		},
		{
			name:       "zero price",
			product:    &domain.Product{Name: "Phone", Price: 0, Quantity: 1, Brand: "Vivo", Category: "Device"},
			wantErrMsg: "product price must be greater than 0",
		},
		{
			name:       "negative price",
			product:    &domain.Product{Name: "Keyboard", Price: -1, Quantity: 1, Brand: "Logitech", Category: "Electronics"},
			wantErrMsg: "product price must be greater than 0",
		},
		{
			name:       "negative quantity",
			product:    &domain.Product{Name: "Keyboard", Price: 600, Quantity: -3, Brand: "Logitech", Category: "Electronics"},
			wantErrMsg: "quantity must be greater than or equal to 0",
		},
		{
			name:       "missing brand",
			product:    &domain.Product{Name: "Keyboard", Price: 600, Quantity: 1, Brand: "", Category: "Electronics"},
			wantErrMsg: "brand name is compulsory for all products",
		},
		{
			name:       "missing category",
			product:    &domain.Product{Name: "Keyboard", Price: 100, Quantity: 1, Brand: "Logitech", Category: ""},
			wantErrMsg: "category name is required",
		},
		{
			name:    "zero quantity is valid",
			product: &domain.Product{Name: "Keyboard", Price: 100, Quantity: 0, Brand: "Logitech", Category: "Electronics"},
		},
		{
			name:            "category not found in repo",
			product:         mockProduct(),
			categoryRepoErr: errors.New("category not found"),
			wantErrMsg:      "category not found",
		},
		{
			name:           "product repository error",
			product:        mockProduct(),
			productRepoErr: errors.New("Failed to add new product"),
			wantErrMsg:     "Failed to add new product",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			categoryRepo := &mockCategoryRepo{
				getByTitle: func(ctx context.Context, title string) (*domain.ProductCategory, error) {
					if tt.categoryRepoErr != nil {
						return nil, tt.categoryRepoErr
					}
					return mockCategory(), nil
				},
			}

			productRepo := &mockProductRepo{
				create: func(ctx context.Context, product *domain.Product) error {
					return tt.productRepoErr
				},
			}

			svc := NewProductService(productRepo, categoryRepo)
			err := svc.CreateProduct(context.Background(), tt.product)

			if tt.wantErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name           string
		product        *domain.Product
		productRepoErr error
		wantErrMsg     string
	}{
		{
			name:    "success",
			product: mockProduct(),
		},
		{
			name:       "missing name",
			product:    &domain.Product{Price: 100, Quantity: 1, Brand: "Logitech"},
			wantErrMsg: "product name is required",
		},
		{
			name:       "negative price",
			product:    &domain.Product{Name: "Keyboard", Quantity: 9, Price: -5, Brand: "Logitech"},
			wantErrMsg: "product price must be greater than 0",
		},
		{
			name:       "negative quantity",
			product:    &domain.Product{Name: "Keyboard", Quantity: -7, Price: 5, Brand: "Logitech"},
			wantErrMsg: "quantity must be greater than or equal to 0",
		},
		{
			name:       "zero price",
			product:    &domain.Product{Name: "Keyboard", Quantity: 9, Price: 0, Brand: "Logitech"},
			wantErrMsg: "product price must be greater than 0",
		},
		{
			name:       "missing brand",
			product:    &domain.Product{Name: "Keyboard", Quantity: 9, Price: 2, Brand: ""},
			wantErrMsg: "brand name is compulsory for all products",
		},
		{
			name:           "product not found",
			product:        mockProduct(),
			productRepoErr: errors.New("product not found"),
			wantErrMsg:     "product not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productRepo := &mockProductRepo{
				update: func(ctx context.Context, id string, product *domain.Product) error {
					return tt.productRepoErr
				},
			}

			svc := NewProductService(productRepo, &mockCategoryRepo{})
			err := svc.UpdateProduct(context.Background(), "123", tt.product)

			if tt.wantErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	tests := []struct {
		name           string
		mockResult     *domain.Product
		productRepoErr error
		wantErrMsg     string
	}{
		{
			name:       "success",
			mockResult: mockProduct(),
		},
		{
			name:           "product not found",
			mockResult:     nil,
			productRepoErr: errors.New("product not found"),
			wantErrMsg:     "product not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productRepo := &mockProductRepo{
				get: func(ctx context.Context, id string) (*domain.Product, error) {
					return tt.mockResult, tt.productRepoErr
				},
			}

			svc := NewProductService(productRepo, &mockCategoryRepo{})
			result, err := svc.GetProduct(context.Background(), "123")

			if tt.wantErrMsg == "" {
				require.NoError(t, err)
				assert.Equal(t, tt.mockResult, result)
			} else {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
				assert.Nil(t, result)
			}
		})
	}
}

func TestListProducts(t *testing.T) {
	tests := []struct {
		name           string
		mockResult     []*domain.Product
		productRepoErr error
		wantErrMsg     string
	}{
		{
			name:       "success",
			mockResult: []*domain.Product{mockProduct(), mockProduct()},
		},
		{
			name:           "repository error",
			mockResult:     nil,
			productRepoErr: errors.New("failed to fetch products"),
			wantErrMsg:     "failed to fetch products",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productRepo := &mockProductRepo{
				list: func(ctx context.Context) ([]*domain.Product, error) {
					return tt.mockResult, tt.productRepoErr
				},
			}

			svc := NewProductService(productRepo, &mockCategoryRepo{})
			result, err := svc.ListProducts(context.Background())

			if tt.wantErrMsg == "" {
				require.NoError(t, err)
				assert.Equal(t, tt.mockResult, result)
			} else {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
				assert.Nil(t, result)
			}
		})
	}
}

func TestGetByCategory(t *testing.T) {
	tests := []struct {
		name            string
		categoryTitle   string
		mockResult      []*domain.Product
		categoryRepoErr error
		productRepoErr  error
		wantErrMsg      string
	}{
		{
			name:          "success",
			categoryTitle: "Electronics",
			mockResult:    []*domain.Product{mockProduct(), mockProduct()},
		},
		{
			// GetByCategory first fetches the category by title, so a missing category should revert back before hitting the product repo
			name:            "category not found",
			categoryTitle:   "Unknown",
			categoryRepoErr: errors.New("category not found"),
			wantErrMsg:      "category not found",
		},
		{
			name:           "product repository error",
			categoryTitle:  "Electronics",
			productRepoErr: errors.New("failed to fetch products"),
			wantErrMsg:     "failed to fetch products",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryRepo := &mockCategoryRepo{
				getByTitle: func(ctx context.Context, title string) (*domain.ProductCategory, error) {
					if tt.categoryRepoErr != nil {
						return nil, tt.categoryRepoErr
					}
					return mockCategory(), nil
				},
			}

			productRepo := &mockProductRepo{
				getByCategory: func(ctx context.Context, id string) ([]*domain.Product, error) {
					return tt.mockResult, tt.productRepoErr
				},
			}

			svc := NewProductService(productRepo, categoryRepo)
			result, err := svc.GetByCategory(context.Background(), tt.categoryTitle)

			if tt.wantErrMsg == "" {
				require.NoError(t, err)
				assert.Equal(t, tt.mockResult, result)
			} else {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
				assert.Nil(t, result)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tests := []struct {
		name           string
		productRepoErr error
		wantErrMsg     string
	}{
		{
			name: "success",
		},
		{
			name:           "product not found",
			productRepoErr: errors.New("product not found"),
			wantErrMsg:     "product not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productRepo := &mockProductRepo{
				delete: func(ctx context.Context, id string) error {
					return tt.productRepoErr
				},
			}

			svc := NewProductService(productRepo, &mockCategoryRepo{})
			err := svc.DeleteProduct(context.Background(), "123")

			if tt.wantErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
			}
		})
	}
}
