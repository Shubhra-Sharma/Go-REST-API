package repository

import (
	"context"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	Get(ctx context.Context, id string) (*domain.Product, error)
	List(ctx context.Context) ([]*domain.Product, error)
	GetByCategory(ctx context.Context, id string) ([]*domain.Product, error)
	Update(ctx context.Context, id string, product *domain.Product) error
	Delete(ctx context.Context, id string) error
}
