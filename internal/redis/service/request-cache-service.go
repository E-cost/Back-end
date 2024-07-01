package service

import (
	"Ecost/internal/utils/variables"
	"Ecost/pkg/logging"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewRedisIpRequestService(
	logger *logging.Logger,
	client *redis.Client,
) RedisIpRequestService {
	return &redisService{
		logger: logger,
		client: client,
	}
}

func (s *redisService) GetRequestCount(ctx context.Context, ip string, path string) (int, error) {
	key := fmt.Sprintf("%s:%s", ip, path)
	result, err := s.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return result, nil
}

func (s *redisService) IncrementRequestCount(ctx context.Context, ip string, path string) error {
	key := fmt.Sprintf("%s:%s", ip, path)
	_, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	_, err = s.client.Expire(ctx, key, variables.CacheTTL).Result()
	if err != nil {
		return err
	}

	return nil
}
