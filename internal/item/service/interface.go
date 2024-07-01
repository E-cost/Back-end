package service

import (
	"Ecost/internal/api/middleware/sort"
	"Ecost/internal/facades"
	"Ecost/internal/utils/types"
	"context"
)

type ItemService interface {
	SetYandexFacade(facade facades.YandexItemFacade)
	GetOneItemSeafood(ctx context.Context, id string) (*types.GetSeafoodType, error)
	GetAllSeafoodSorted(ctx context.Context, sortOptions sort.Options) ([]*types.GetSeafoodType, error)
	GetAllSeafood(ctx context.Context) ([]*types.GetSeafoodType, error)
	GetAllIdsSeafood(ctx context.Context) ([]string, error)
}
