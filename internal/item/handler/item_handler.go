package itemsHandler

import (
	"Ecost/internal/api/handlers"
	"Ecost/internal/api/middleware/errors"
	"Ecost/internal/api/middleware/limit"
	"Ecost/internal/api/middleware/sort"
	"Ecost/internal/apperror"
	"Ecost/internal/facades"
	service1 "Ecost/internal/item/service"
	"Ecost/internal/item/storage"
	service2 "Ecost/internal/redis/service"
	"Ecost/internal/utils/types"
	"Ecost/pkg/logging"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	GETItemSeafoodURL  = "/items/seafood/:uuid"
	GETItemsSeafoodURL = "/items/seafood"
	GETItemsAll        = "/items/all"
)

var sortOptions sort.Options

type handler struct {
	logger     *logging.Logger
	service    service1.ItemService
	repository storage.ItemsRepository
}

func NewItemHandler(
	repository storage.ItemsRepository,
	logger *logging.Logger,
	redisItemsSvc service2.RedisItemService,
	yandexFacade facades.YandexItemFacade,
) handlers.ItemHandler {
	itemService := service1.NewItemService(logger, repository, redisItemsSvc, yandexFacade)
	return &handler{
		logger:     logger,
		service:    itemService,
		repository: repository,
	}
}

func (h *handler) SaveItem(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, GETItemSeafoodURL, limit.ReadMiddleware(errors.Middleware(h.GetItemSeafood)))
	router.HandlerFunc(http.MethodGet, GETItemsSeafoodURL, limit.ReadMiddleware(sort.Middleware(errors.Middleware(h.GetListOfSeafood), "created_at", sort.ASC)))
	router.HandlerFunc(http.MethodGet, GETItemsAll, limit.ReadMiddleware(errors.Middleware(h.GetAllItems)))
}

func (h *handler) GetItemSeafood(w http.ResponseWriter, r *http.Request) error {
	uuid := httprouter.ParamsFromContext(r.Context()).ByName("uuid")
	if uuid == "" {
		return apperror.BadRequest("id has not been provided", "")
	}

	one, err := h.service.GetOneItemSeafood(r.Context(), uuid)
	if err != nil {
		return apperror.NotFound("item not found", err.Error())
	}

	responseBytes, err := json.Marshal(one)
	if err != nil {
		return apperror.InternalServerError("Failed to marshal item's data", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

	return nil
}

func (h *handler) GetListOfSeafood(w http.ResponseWriter, r *http.Request) error {
	if options, ok := r.Context().Value(sort.OptionsContextKey).(sort.Options); ok {
		sortOptions = options
	}

	all, err := h.service.GetAllSeafoodSorted(r.Context(), sortOptions)
	if err != nil {
		w.WriteHeader(400)
		return err
	}

	allBytes, err := json.Marshal(all)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)

	return nil
}

func (h *handler) GetAllItems(w http.ResponseWriter, r *http.Request) error {
	seafood, err := h.service.GetAllSeafood(r.Context())
	if err != nil {
		w.WriteHeader(400)
		return err
	}

	response := &types.AllItemsResponse{
		Seafood: seafood,
	}

	allBytes, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)

	return nil
}
