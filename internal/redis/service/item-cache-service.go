package service

import (
	"Ecost/internal/utils/types"
	"Ecost/internal/utils/variables"
	"Ecost/pkg/logging"
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

func NewRedisItemsService(
	logger *logging.Logger,
	client *redis.Client,
) RedisItemService {
	return &redisService{
		logger: logger,
		client: client,
	}
}

func (s *redisService) GetAllItemsCache(ctx context.Context, key string) ([]*types.GetSeafoodType, error) {
	data, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var items []*types.GetSeafoodType
	err = json.Unmarshal([]byte(data), &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *redisService) GetSeafoodItemsCache(ctx context.Context, key string) ([]*types.GetSeafoodType, error) {
	data, err := s.client.Get(ctx, variables.AllSeafoodItems+key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var items []*types.GetSeafoodType
	err = json.Unmarshal([]byte(data), &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *redisService) GetSeafoodItem(ctx context.Context, key string) (*types.GetSeafoodType, error) {
	data, err := s.client.Get(ctx, variables.ItemKey+key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var item types.GetSeafoodType
	err = json.Unmarshal([]byte(data), &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *redisService) SetAllItemsCache(ctx context.Context, items []*types.GetSeafoodType) error {
	key := variables.AllItemsKey
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}

	err = s.client.Set(ctx, key, data, variables.CacheTTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *redisService) SetSeafoodItemsCache(ctx context.Context, items []*types.GetSeafoodType) error {
	key := variables.AllSeafoodItems
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}

	err = s.client.Set(ctx, key, data, variables.CacheTTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *redisService) SetSeafoodItem(ctx context.Context, item *types.GetSeafoodType) error {
	key := variables.ItemKey + item.ID
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	err = s.client.Set(ctx, key, data, variables.CacheTTL).Err()
	if err != nil {
		return err
	}

	return nil
}
