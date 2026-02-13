package categories

import (
	"context"

	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
)

type Service interface {
	GetCategories(ctx context.Context) ([]sqlc.Category, error)
}
