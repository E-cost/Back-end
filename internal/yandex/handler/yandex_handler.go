package yandexHandler

import (
	"Ecost/internal/api/handlers"
	"Ecost/internal/api/middleware/errors"
	"Ecost/internal/api/middleware/limit"
	"Ecost/internal/apperror"
	"Ecost/internal/config"
	"Ecost/internal/yandex/service"
	"Ecost/pkg/logging"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/julienschmidt/httprouter"
)

const (
	GETPhotoById = "/photos/:uuid"
	GETPriceList = "/price-list/download"
)

type handler struct {
	logger  *logging.Logger
	service service.YandexService
}

func NewYandexHandler(
	logger *logging.Logger,
	client *s3.PresignClient,
	config *config.Config,
) handlers.YandexHandler {

	yandexService := service.NewYandexService(logger, client, config)
	return &handler{
		logger:  logger,
		service: yandexService,
	}
}

func (h *handler) SaveYandexStorage(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, GETPriceList, limit.ReadMiddleware(errors.Middleware(h.GetPriceList)))
	router.HandlerFunc(http.MethodGet, GETPhotoById, limit.ReadMiddleware(errors.Middleware(h.GetPhotosById)))
}

func (h *handler) GetPriceList(w http.ResponseWriter, r *http.Request) error {
	data, contentType, err := h.service.DownloadPriceList(r.Context())
	if err != nil {
		h.logger.Errorf("Failed to get price-list from storage: %v", err)
		apperror.InternalServerError("Failed to get price-list", err.Error())
		return nil
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename=price-list.xlsx")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		h.logger.Errorf("Failed to write data to response: %v", err)
		apperror.InternalServerError("Failed to write data to response", err.Error())
	}

	return nil
}

func (h *handler) GetPhotosById(w http.ResponseWriter, r *http.Request) error {
	uuid := httprouter.ParamsFromContext(r.Context()).ByName("uuid")
	if uuid == "" {
		return apperror.BadRequest("id has not been provided", "")
	}

	links, err := h.service.GetItemPhotos(context.Background(), uuid)
	if err != nil {
		h.logger.Errorf("Failed to get photos by id of the product from storage: %v", err.Error())
		return apperror.InternalServerError("Failed to get photo", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(links); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err.Error())
		return apperror.InternalServerError("Failed to encode response", err.Error())
	}

	return nil
}
