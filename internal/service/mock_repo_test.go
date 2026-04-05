package service

import (
	"context"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
)

// Mocking ProductRepository Interface
type mockProductRepo struct {
	create           func(ctx context.Context, product *domain.Product) error
	get              func(ctx context.Context, id string) (*domain.Product, error)
	list             func(ctx context.Context) ([]*domain.Product, error)
	getByCategory    func(ctx context.Context, id string) ([]*domain.Product, error)
	update           func(ctx context.Context, id string, product *domain.Product) error
	delete           func(ctx context.Context, id string) error
	deleteByCategory func(ctx context.Context, id string) error
}

func (m *mockProductRepo) Create(ctx context.Context, product *domain.Product) error {
	return m.create(ctx, product)
}

func (m *mockProductRepo) Get(ctx context.Context, id string) (*domain.Product, error) {
	return m.get(ctx, id)
}

func (m *mockProductRepo) List(ctx context.Context) ([]*domain.Product, error) {
	return m.list(ctx)
}

func (m *mockProductRepo) GetByCategory(ctx context.Context, id string) ([]*domain.Product, error) {
	return m.getByCategory(ctx, id)
}

func (m *mockProductRepo) Update(ctx context.Context, id string, product *domain.Product) error {
	return m.update(ctx, id, product)
}

func (m *mockProductRepo) Delete(ctx context.Context, id string) error {
	return m.delete(ctx, id)
}

func (m *mockProductRepo) DeleteByCategoryID(ctx context.Context, id string) error {
	return m.deleteByCategory(ctx, id)
}

// Mocking ProductCategoryRepository Interface
type mockCategoryRepo struct {
	create     func(ctx context.Context, category *domain.ProductCategory) (*domain.ProductCategory, error)
	list       func(ctx context.Context) ([]*domain.ProductCategory, error)
	getByTitle func(ctx context.Context, title string) (*domain.ProductCategory, error)
	update     func(ctx context.Context, id string, category *domain.ProductCategory) error
	delete     func(ctx context.Context, id string) error
}

func (m *mockCategoryRepo) Create(ctx context.Context, category *domain.ProductCategory) (*domain.ProductCategory, error) {
	return m.create(ctx, category)
}

func (m *mockCategoryRepo) List(ctx context.Context) ([]*domain.ProductCategory, error) {
	return m.list(ctx)
}

func (m *mockCategoryRepo) GetByTitle(ctx context.Context, title string) (*domain.ProductCategory, error) {
	return m.getByTitle(ctx, title)
}

func (m *mockCategoryRepo) Update(ctx context.Context, id string, category *domain.ProductCategory) error {
	return m.update(ctx, id, category)
}

func (m *mockCategoryRepo) Delete(ctx context.Context, id string) error {
	return m.delete(ctx, id)
}

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
