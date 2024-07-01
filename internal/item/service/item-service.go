package service

import (
	"Ecost/internal/api/middleware/sort"
	"Ecost/internal/apperror"
	"Ecost/internal/facades"
	"Ecost/internal/item/service/helpers"
	"Ecost/internal/item/storage"
	"Ecost/internal/redis/service"
	sort2 "Ecost/internal/utils/sort"
	"Ecost/internal/utils/types"
	"Ecost/internal/utils/variables"
	"Ecost/pkg/logging"
	"context"
)

type itemService struct {
	logger       *logging.Logger
	repository   storage.ItemsRepository
	redisItems   service.RedisItemService
	yandexFacade facades.YandexItemFacade
}

func NewItemService(
	logger *logging.Logger,
	repository storage.ItemsRepository,
	redisItems service.RedisItemService,
	yandexFacade facades.YandexItemFacade,
) ItemService {
	return &itemService{
		logger:       logger,
		repository:   repository,
		redisItems:   redisItems,
		yandexFacade: yandexFacade,
	}
}

func (s *itemService) SetYandexFacade(facade facades.YandexItemFacade) {
	s.yandexFacade = facade
}

func (s *itemService) GetOneItemSeafood(ctx context.Context, id string) (*types.GetSeafoodType, error) {
	// Check if the item is in cache
	cachedItem, err := s.redisItems.GetSeafoodItem(ctx, id)
	if err == nil && cachedItem != nil {
		return cachedItem, nil
	}

	// The item not found in cache, get it from the database
	item, err := s.repository.FindOneSeafood(ctx, id)
	if err != nil {
		return nil, apperror.InternalServerError("failed to get the item by UUID: %v", err.Error())
	}

	if item.ID == "" {
		return nil, apperror.NotFound("no such item: %v", "no such item")
	}

	itemType, err := helpers.ConvertToSeafoodType(item)
	if err != nil {
		return nil, apperror.InternalServerError("failed to convert an item due error: %v", err.Error())
	}

	// Cache the item with TTL of 6 hours
	err = s.redisItems.SetSeafoodItem(ctx, itemType)
	if err != nil {
		return nil, apperror.InternalServerError("failed to get the item by UUID: %v", err.Error())
	}

	return itemType, nil
}

func (s *itemService) GetAllSeafoodSorted(ctx context.Context, sortOptions sort.Options) ([]*types.GetSeafoodType, error) {
	// Check if all items are in cache
	cachedItems, err := s.redisItems.GetSeafoodItemsCache(ctx, variables.AllSeafoodItems)
	if err == nil && cachedItems != nil {
		return cachedItems, nil
	}

	options := sort2.NewSortOptions(sortOptions.Field, sortOptions.Order)
	all, err := s.repository.FindAllSeafood(ctx, options)
	if err != nil {
		return nil, apperror.InternalServerError("failed to get all items due to error: %v", err.Error())
	}

	var ids []string
	for _, item := range all {
		ids = append(ids, item.ID)
	}

	urls, err := s.yandexFacade.GetPreviewPhotos(ctx, ids)
	if err != nil {
		return nil, apperror.InternalServerError("failed to call yandex GetAllItemsIds due an error: %v", err.Error())
	}

	itemsType, err := helpers.ConvertToSeafoodSliceType(all, urls)
	if err != nil {
		return nil, apperror.InternalServerError("failed to convert items due to error: %v", err.Error())
	}

	// Cache the item with TTL of 6 hours
	err = s.redisItems.SetSeafoodItemsCache(ctx, itemsType)
	if err != nil {
		s.logger.Warnf("Error setting items in cache: %v", err)
	}

	return itemsType, nil
}

func (s *itemService) GetAllSeafood(ctx context.Context) ([]*types.GetSeafoodType, error) {
	// Check if all items are in cache
	cachedItems, err := s.redisItems.GetAllItemsCache(ctx, variables.AllItemsKey)
	if err == nil && cachedItems != nil {
		return cachedItems, nil
	}

	// Items not found in cache, get them from the database
	all, err := s.repository.GetAllSeafood(ctx)
	if err != nil {
		return nil, apperror.InternalServerError("failed to get all items due to error: %v", err.Error())
	}

	var ids []string
	for _, item := range all {
		ids = append(ids, item.ID)
	}

	urls, err := s.yandexFacade.GetPreviewPhotos(ctx, ids)
	if err != nil {
		return nil, apperror.InternalServerError("failed to call yandex GetAllItemsIds due an error: %v", err.Error())
	}

	itemsType, err := helpers.ConvertToSeafoodSliceType(all, urls)
	if err != nil {
		return nil, apperror.InternalServerError("failed to convert items due to error: %v", err.Error())
	}

	// Cache the items with TTL of 6 hours
	err = s.redisItems.SetAllItemsCache(ctx, itemsType)
	if err != nil {
		s.logger.Warnf("Error setting items in cache: %v", err)
	}

	return itemsType, nil
}

func (s *itemService) GetAllIdsSeafood(ctx context.Context) ([]string, error) {
	ids, err := s.repository.GetAllIdsSeafood(ctx)
	if err != nil {
		return nil, apperror.InternalServerError("failed to get all ids due to error: %v", err.Error())
	}

	return ids, nil
}
