package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
)

// Mocking product and Category structs for testcases
func mockProduct() *domain.Product {
	return &domain.Product{
		Name:     "iPhone 15",
		Price:    999,
		Quantity: 10,
		Brand:    "Apple",
		Category: "Electronics",
	}
}

func mockCategory() *domain.ProductCategory {
	return &domain.ProductCategory{
		ID:    "234",
		Title: "Electronics",
	}
}

// Table Driven unit test for createProduct Service
func TestCreateProduct(t *testing.T) {
	testcases := []struct {
		name            string
		product         *domain.Product
		categoryRepoErr error
		productRepoErr  error
		want            string
	}{
		{ // default result
			name:    "success",
			product: mockProduct(),
			want:    "",
		},
		{ // validation logic tests
			name:    "missing name",
			product: &domain.Product{Price: 100, Quantity: 1, Brand: "Vivo", Category: "Device"},
			want:    "product name is required",
		},
		{
			name:    "zero price",
			product: &domain.Product{Name: "Phone", Price: 0, Quantity: 1, Brand: "Vivo", Category: "Device"},
			want:    "product price must be greater than 0",
		},
		{
			name:    "negative price",
			product: &domain.Product{Name: "Keyboard", Price: -1, Quantity: 1, Brand: "Logitech", Category: "Electronics"},
			want:    "product price must be greater than 0",
		},
		{
			name:    "negative quantity",
			product: &domain.Product{Name: "Keyboard", Price: 600, Quantity: -3, Brand: "Logitech", Category: "Electronics"},
			want:    "quantity must be greater than or equal to 0",
		},
		{
			name:    "missing brand",
			product: &domain.Product{Name: "Keyboard", Price: 600, Quantity: 1, Brand: "", Category: "Electronics"},
			want:    "brand name is compulsory for all products",
		},
		{
			name:    "missing category",
			product: &domain.Product{Name: "Keyboard", Price: 100, Quantity: 1, Brand: "Logitech", Category: ""},
			want:    "category name is required",
		},
		{
			name:    "zero quantity is valid",
			product: &domain.Product{Name: "Keyboard", Price: 100, Quantity: 0, Brand: "Logitech", Category: "Electronics"},
			want:    "",
		},
		{ // To check if service is correctly propagating repository errors
			name:            "category not found in repo",
			product:         mockProduct(),
			categoryRepoErr: errors.New("category not found"),
			want:            "category not found",
		},
		{ // To check if service is correctly propagating repository errors
			name:           "Product repository error",
			product:        mockProduct(),
			productRepoErr: errors.New("Failed to add new product"),
			want:           "Failed to add new product",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			// Mock categoryRepo
			categoryRepo := &mockCategoryRepo{
				getByTitle: func(ctx context.Context, title string) (*domain.ProductCategory, error) {
					if tt.categoryRepoErr != nil {
						return nil, tt.categoryRepoErr
					}
					return mockCategory(), nil
				},
			}

			// mock productRepo
			productRepo := &mockProductRepo{
				create: func(ctx context.Context, product *domain.Product) error {
					if tt.productRepoErr != nil {
						return tt.productRepoErr
					}
					return nil
				},
			}

			// creating new service and calling the create function for testing
			service := NewProductService(productRepo, categoryRepo)
			err := service.CreateProduct(context.Background(), tt.product)

			if tt.want == "" && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name           string
		product        *domain.Product
		productRepoErr error
		want           string
	}{
		{ // default
			name:    "success",
			product: mockProduct(),
			want:    "",
		},
		{ // testing validation logic
			name:    "missing name",
			product: &domain.Product{Price: 100, Quantity: 1, Brand: "Logitech"},
			want:    "product name is required",
		},
		{
			name:    "negative price",
			product: &domain.Product{Name: "Keyboard", Quantity: 9, Price: -5, Brand: "Logitech"},
			want:    "product price must be greater than 0",
		},
		{
			name:    "negative quantity",
			product: &domain.Product{Name: "Keyboard", Quantity: -7, Price: 5, Brand: "Logitech"},
			want:    "quantity must be greater than or equal to 0",
		},
		{
			name:    "zero price",
			product: &domain.Product{Name: "Keyboard", Quantity: 9, Price: 0, Brand: "Logitech"},
			want:    "product price must be greater than 0",
		},
		{
			name:    "missing brand",
			product: &domain.Product{Name: "Keyboard", Quantity: 9, Price: 2, Brand: ""},
			want:    "brand name is compulsory for all products",
		},
		{ // testing if ProductService correctly propagates repository errors
			name:           "product not found",
			product:        mockProduct(),
			productRepoErr: errors.New("product not found"),
			want:           "product not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mocking productRepo
			productRepo := &mockProductRepo{
				update: func(ctx context.Context, id string, product *domain.Product) error {
					return tt.productRepoErr
				},
			}

			service := NewProductService(productRepo, &mockCategoryRepo{}) // passing empty mockCategory struct since there is no role of category Repo in UpdateProduct function
			err := service.UpdateProduct(context.Background(), "123", tt.product)

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
