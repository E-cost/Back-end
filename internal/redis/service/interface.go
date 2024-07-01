package service

import (
	"Ecost/internal/utils/types"
	"Ecost/pkg/logging"
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisIpRequestService interface {
	IncrementRequestCount(ctx context.Context, ip string, path string) error
	GetRequestCount(ctx context.Context, ip string, path string) (int, error)
}

type RedisItemService interface {
	GetAllItemsCache(ctx context.Context, key string) ([]*types.GetSeafoodType, error)
	GetSeafoodItemsCache(ctx context.Context, key string) ([]*types.GetSeafoodType, error)
	GetSeafoodItem(ctx context.Context, key string) (*types.GetSeafoodType, error)
	SetAllItemsCache(ctx context.Context, items []*types.GetSeafoodType) error
	SetSeafoodItemsCache(ctx context.Context, items []*types.GetSeafoodType) error
	SetSeafoodItem(ctx context.Context, item *types.GetSeafoodType) error
}

type redisService struct {
	logger *logging.Logger
	client *redis.Client
}
