package service

import (
	"context"
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
					return tt.productRepoErr
				},
			}

			// creating new service and calling the create function for testing
			svc := NewProductService(productRepo, categoryRepo)
			err := svc.CreateProduct(context.Background(), tt.product)

			if tt.want == "" && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}
