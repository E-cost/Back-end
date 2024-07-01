package service

import (
	"Ecost/internal/apperror"
	"Ecost/internal/facades"
	service2 "Ecost/internal/order/service"
	"Ecost/internal/yandex/service"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

type yandexOrderFacade struct {
	yandexService service.YandexService
	orderService  service2.OrderService
}

func NewYandexOrderFacade(yandexService service.YandexService, orderService service2.OrderService) facades.YandexOrderFacade {
	return &yandexOrderFacade{
		yandexService: yandexService,
		orderService:  orderService,
	}
}

func (f *yandexOrderFacade) PutOrder(ctx context.Context, key string, content []byte) error {
	req, err := f.yandexService.PutObject(ctx, key, content)
	if err != nil {
		return apperror.InternalServerError("Failed to call yandex PutObject method due an error: %v", err.Error())
	}

	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	putReq, err := http.NewRequest(http.MethodPut, req.URL, bytes.NewReader(content))
	if err != nil {
		return apperror.InternalServerError("Failed to create PUT request: %v", err.Error())
	}

	putReq.Header.Set("Content-Type", "application/octet-stream")
	putReq.ContentLength = int64(len(content))

	resp, err := httpClient.Do(putReq)
	if err != nil {
		return apperror.InternalServerError("Failed to execute PUT request: %v", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apperror.InternalServerError("PUT request failed with status: %d", resp.Status)
	}

	fmt.Println("File uploaded successfully.") // TODO: to add a logger

	return nil
}

func (f *yandexOrderFacade) DeleteOrder(ctx context.Context, keys []string) error {
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	for _, key := range keys {
		req, err := f.yandexService.DeleteObject(ctx, key)
		if err != nil {
			return apperror.InternalServerError("Failed to delete order from Yandex Object Storage: %v", err.Error())
		}

		delReq, err := http.NewRequest(http.MethodDelete, req.URL, nil)
		if err != nil {
			return apperror.InternalServerError("Failed to create DELETE request: %v", err.Error())
		}

		resp, err := httpClient.Do(delReq)
		if err != nil {
			return apperror.InternalServerError("Failed to execute DELETE request: %v", err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			return apperror.InternalServerError("DELETE request failed with status: %d", resp.Status)
		}

		fmt.Println("File(s) deleted successfully.") // TODO: to add a logger
	}

	return nil
}
