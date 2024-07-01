package service

import (
	"Ecost/internal/apperror"
	"Ecost/internal/facades"
	service2 "Ecost/internal/item/service"
	"Ecost/internal/yandex/service"
	"context"
)

type yandexItemFacade struct {
	yandexService service.YandexService
	itemService   service2.ItemService
}

func NewYandexItemFacade(yandexService service.YandexService, itemService service2.ItemService) facades.YandexItemFacade {
	return &yandexItemFacade{
		yandexService: yandexService,
		itemService:   itemService,
	}
}

func (f *yandexItemFacade) GetPreviewPhotos(ctx context.Context, ids []string) (map[string]string, error) {
	photoURLs := make(map[string]string)
	for _, id := range ids {
		url, err := f.yandexService.GetPreviewPhoto(ctx, id)
		if err != nil {
			return nil, apperror.InternalServerError("Failed to get preview photo for id: %v", err.Error())
		}
		photoURLs[id] = url
	}
	return photoURLs, nil
}
