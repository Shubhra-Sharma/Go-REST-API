package repository

import (
	"context"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
)

type ProductCategoryRepository interface {
	Create(ctx context.Context, product *domain.ProductCategory) error
	Get(ctx context.Context, id string) (*domain.ProductCategory, error)
	List(ctx context.Context) ([]*domain.ProductCategory, error)
	Update(ctx context.Context, id string, product *domain.ProductCategory) error
	Delete(ctx context.Context, id string) error
}
