package service

import (
	"Ecost/internal/facades"
	service2 "Ecost/internal/item/service"
	"Ecost/internal/order/model"
	"Ecost/internal/order/service"
	"Ecost/internal/utils/types"
	"context"
	"errors"
)

type orderItemFacade struct {
	orderService service.OrderService
	itemService  service2.ItemService
}

func NewOrderItemFacade(orderService service.OrderService, itemService service2.ItemService) facades.OrderItemFacade {
	return &orderItemFacade{
		orderService: orderService,
		itemService:  itemService,
	}
}

func (f *orderItemFacade) VerifyItemById(ctx context.Context, id string, category string, quantity int) (*model.Item, error) {
	switch category {
	case "seafood":
		item, err := f.itemService.GetOneItemSeafood(ctx, id)
		if err != nil {
			return nil, err
		}

		if item.ID == id {
			return &model.Item{
				ID:        item.ID,
				Article:   item.Article,
				Category:  item.Category,
				Product:   item.Product,
				Name:      item.Name,
				NetWeight: item.NetWeight,
				Quantity:  quantity,
				PriceBYN:  *item.PriceBYN,
			}, nil
		}
	default:
		return nil, errors.New("no such item")
	}
	return nil, nil
}

func (f *orderItemFacade) VerifyItems(ctx context.Context, items []types.EnterOrderItem) ([]facades.ItemResult, error) {
	results := make([]facades.ItemResult, len(items))

	for i, item := range items {
		validItem, err := f.VerifyItemById(ctx, item.ID, item.Category, item.Quantity)
		if err != nil {
			return nil, err
		}
		if validItem != nil {
			results[i] = facades.ItemResult{Valid: true, Item: *validItem}
		} else {
			results[i] = facades.ItemResult{Valid: false, Item: model.Item{}}
		}
	}
	return results, nil
}
