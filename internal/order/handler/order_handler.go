package ordersHandler

import (
	"Ecost/internal/api/handlers"
	"Ecost/internal/api/middleware/dns"
	"Ecost/internal/api/middleware/errors"
	"Ecost/internal/api/middleware/limit"
	"Ecost/internal/apperror"
	"Ecost/internal/facades"
	"Ecost/internal/order/dto"
	"Ecost/internal/order/service"
	"Ecost/internal/order/storage"
	"Ecost/internal/utils/ip"
	"Ecost/pkg/logging"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	POSTOrderCreateURL = "/orders/create"
)

type handler struct {
	logger       *logging.Logger
	ipMiddleware *dns.Middleware
	service      service.OrderService
	repository   storage.OrdersRepository
}

func NewOrderHandler(
	logger *logging.Logger,
	ipMiddleware *dns.Middleware,
	contactFacade facades.ContactOrderFacade,
	itemFacade facades.OrderItemFacade,
	yandexFacade facades.YandexOrderFacade,
	repository storage.OrdersRepository,
) handlers.OrderHandler {

	orderService := service.NewOrderService(logger, repository, contactFacade, itemFacade, yandexFacade)
	return &handler{
		logger:       logger,
		ipMiddleware: ipMiddleware,
		service:      orderService,
		repository:   repository,
	}
}

func (h *handler) SaveOrder(router *httprouter.Router) {
	router.HandlerFunc(
		http.MethodPost,
		POSTOrderCreateURL,
		limit.WriteMiddleware(h.ipMiddleware.IpMiddleware(errors.Middleware(h.CreateOrder))))
}

func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	var createOrderDto dto.CreateOrderDto

	IPOutput, err := ip.ReadUserIP(r)
	if err != nil {
		return apperror.InternalServerError("Failed to get ip from request due an error:", err.Error())
	}

	if err := json.NewDecoder(r.Body).Decode(&createOrderDto); err != nil {
		return apperror.BadRequest("Invalid request body", err.Error())
	}

	output, pdfContent, err := h.service.CreateOrder(r.Context(), createOrderDto, *IPOutput)
	if err != nil {
		return apperror.BadRequest("Failed to create an order", err.Error())
	}

	w.Header().Set("Content-Disposition", "attachment; filename=order.pdf")
	w.Header().Set("Content-Type", "application/pdf")

	if _, err := w.Write(pdfContent); err != nil {
		return apperror.InternalServerError("Failed to write PDF content to response", err.Error())
	}

	response := struct {
		EmailId     string `json:"email_id"`
		SecreteCode string `json:"secrete_code"`
		ID          string `json:"id"`
	}{
		EmailId:     output.EmailId,
		SecreteCode: output.SecretCode,
		ID:          output.OrderId,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return apperror.InternalServerError("Failed to encode response", err.Error())
	}

	return nil
}
