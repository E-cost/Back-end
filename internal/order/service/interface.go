package service

import (
	"Ecost/internal/facades"
	"Ecost/internal/order/dto"
	"Ecost/internal/utils/ip"
	"Ecost/internal/utils/types"
	"context"
)

type OrderService interface {
	SetContactFacade(facade facades.ContactOrderFacade)
	SetItemFacade(facade facades.OrderItemFacade)
	SetYandexFacade(facade facades.YandexOrderFacade)
	CreateOrder(ctx context.Context, dto dto.CreateOrderDto, ip ip.IPOutput) (*types.OrderOutput, []byte, error)
	DeleteExpiredOrders(ctx context.Context) error
	initCronJobs()
	StopCronJobs()
}
