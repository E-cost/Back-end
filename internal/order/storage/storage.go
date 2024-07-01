package storage

import (
	"Ecost/internal/order/model"
	"context"
)

type OrdersRepository interface {
	Create(ctx context.Context, order *model.Order) (string, string, float64, error)
	DeleteExpired(ctx context.Context) ([]string, error)
}
