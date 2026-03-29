package repository

import (
	"context"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
)

type ProductCategoryRepository interface {
	Create(ctx context.Context, category *domain.ProductCategory) (*domain.ProductCategory, error)
	List(ctx context.Context) ([]*domain.ProductCategory, error)
	GetByTitle(ctx context.Context, title string) (*domain.ProductCategory, error)
	Update(ctx context.Context, id string, category *domain.ProductCategory) error
	Delete(ctx context.Context, id string) error
}
