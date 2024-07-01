package facades

import (
	"Ecost/internal/order/dto"
	"Ecost/internal/order/model"
	"Ecost/internal/utils/ip"
	"Ecost/internal/utils/types"
	"context"
)

type ContactOrderFacade interface {
	CreateContactOrder(ctx context.Context, dto dto.CreateOrderDto, ip ip.IPOutput) (types.ContactConfirmationInfo, error)
	GetContactOrder(ctx context.Context, id string) (types.ContactInfo, error)
}

type OrderItemFacade interface {
	VerifyItemById(ctx context.Context, id string, category string, quantity int) (*model.Item, error)
	VerifyItems(ctx context.Context, items []types.EnterOrderItem) ([]ItemResult, error)
}

type YandexItemFacade interface {
	GetPreviewPhotos(ctx context.Context, ids []string) (map[string]string, error)
}

type YandexOrderFacade interface {
	PutOrder(ctx context.Context, key string, content []byte) error
	DeleteOrder(ctx context.Context, keys []string) error
}

type ItemResult struct {
	Valid bool
	Item  model.Item
}
