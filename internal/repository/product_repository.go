package repository
import(
	"context"
)

type ProductRepository interface {
    Create(ctx context.Context, product *domain.Product) error
    Get(ctx context.Context, id string) (*domain.Product, error)
    List(ctx context.Context) ([]*domain.Product, error)
    Update(ctx context.Context, id string, product *domain.Product) error
    Delete(ctx context.Context, id string) error
}