package service

import (
	"Ecost/internal/apperror"
	"Ecost/internal/facades"
	"Ecost/internal/order/dto"
	"Ecost/internal/order/model"
	"Ecost/internal/order/service/helpers"
	"Ecost/internal/order/storage"
	"Ecost/internal/utils/ip"
	"Ecost/internal/utils/types"
	"Ecost/pkg/logging"
	"context"
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

type orderService struct {
	logger        *logging.Logger
	cron          *cron.Cron
	repository    storage.OrdersRepository
	contactFacade facades.ContactOrderFacade
	itemFacade    facades.OrderItemFacade
	yandexFacade  facades.YandexOrderFacade
}

var (
	instance *orderService
	once     sync.Once
)

func NewOrderService(
	logger *logging.Logger,
	repository storage.OrdersRepository,
	ctsFacade facades.ContactOrderFacade,
	itsFacade facades.OrderItemFacade,
	ydxFacade facades.YandexOrderFacade,
) OrderService {
	once.Do(func() {
		instance = &orderService{
			logger:        logger,
			cron:          cron.New(),
			repository:    repository,
			contactFacade: ctsFacade,
			itemFacade:    itsFacade,
			yandexFacade:  ydxFacade,
		}

		instance.initCronJobs()
	})
	return instance
}

func (s *orderService) SetContactFacade(facade facades.ContactOrderFacade) {
	s.contactFacade = facade
}

func (s *orderService) SetItemFacade(facade facades.OrderItemFacade) {
	s.itemFacade = facade
}

func (s *orderService) SetYandexFacade(facade facades.YandexOrderFacade) {
	s.yandexFacade = facade
}

func (s *orderService) CreateOrder(ctx context.Context, dto dto.CreateOrderDto, ip ip.IPOutput) (*types.OrderOutput, []byte, error) {
	if err := dto.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validation error: %v", err)
	}

	items := make([]types.EnterOrderItem, len(dto.Items))
	for i, item := range dto.Items {
		items[i] = types.EnterOrderItem{ID: item.ID, Category: item.Category, Quantity: item.Quantity}
	}

	results, err := s.itemFacade.VerifyItems(ctx, items)
	if err != nil {
		return nil, nil, apperror.InternalServerError("Failed to verify items", err.Error())
	}

	validItems := helpers.GetValidItems(results)
	if len(validItems) == 0 {
		return nil, nil, apperror.BadRequest("no such item or items", "")
	}

	orderInfo, err := s.contactFacade.CreateContactOrder(ctx, dto, ip)
	if err != nil {
		return nil, nil, apperror.InternalServerError("Failed to create contact for order", err.Error())
	}

	newOrder := &model.Order{
		ContactId:  orderInfo.ContactId,
		Items:      validItems,
		IpAddress:  ip.RealIP,
		Port:       ip.Port,
		ProxyChain: ip.ProxyChain,
	}

	orderId, shortId, totalSum, err := s.repository.Create(ctx, newOrder)
	if err != nil {
		return nil, nil, apperror.InternalServerError("Failed to create an order", err.Error())
	}

	output := &types.OrderOutput{
		EmailId:    orderInfo.EmailId,
		SecretCode: orderInfo.SecretCode,
		OrderId:    orderId,
	}

	contactInfo, err := s.contactFacade.GetContactOrder(ctx, orderInfo.ContactId)
	if err != nil {
		return nil, nil, apperror.InternalServerError("Failed to get contact info", err.Error())
	}

	pdfContent, err := helpers.GetPDF(newOrder, contactInfo, orderId, shortId, totalSum)
	if err != nil {
		return nil, nil, apperror.InternalServerError("Failed to generate PDF", err.Error())
	}

	key := fmt.Sprintf("orders/%s.pdf", shortId)

	err = s.yandexFacade.PutOrder(ctx, key, pdfContent)
	if err != nil {
		return nil, nil, apperror.InternalServerError("Failed to upload PDF to Yandex Object Storage", err.Error())
	}

	return output, pdfContent, nil
}

func (s *orderService) DeleteExpiredOrders(ctx context.Context) error {
	expiredOrders, err := s.repository.DeleteExpired(ctx)

	var keys []string
	for _, order := range expiredOrders {
		key := fmt.Sprintf("orders/%s.pdf", order)
		keys = append(keys, key)
	}

	if len(keys) > 0 {
		err = s.yandexFacade.DeleteOrder(ctx, keys)
		if err != nil {
			return apperror.InternalServerError("Failed to delete orders from Yandex Object Storage", err.Error())
		}
	}

	if err != nil {
		return apperror.InternalServerError("Failed to delete expired orders", err.Error())
	}

	return nil
}

func (s *orderService) initCronJobs() {
	_, err := s.cron.AddFunc("@weekly", func() {
		ctx := context.Background()
		if err := s.DeleteExpiredOrders(ctx); err != nil {
			s.logger.Errorf("Error deleting expired orders: %v", err)
		}
	})
	if err != nil {
		s.logger.Fatalf("Error scheduling cron job: %v", err)
	}

	s.cron.Start()
}

func (s *orderService) StopCronJobs() {
	s.cron.Stop()
}
