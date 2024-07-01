package storage

import (
	"Ecost/internal/item/model"
	"Ecost/internal/utils/sort"
	"context"
)

type ItemsRepository interface {
	FindOneSeafood(ctx context.Context, id string) (model.Seafood, error)
	FindAllSeafood(ctx context.Context, options sort.SortOptions) ([]model.Seafood, error)
	GetAllSeafood(ctx context.Context) ([]model.Seafood, error)
	GetAllIdsSeafood(ctx context.Context) ([]string, error)
}
