package service

import (
	"context"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

type YandexService interface {
	DownloadPriceList(ctx context.Context) ([]byte, string, error)
	GetItemPhotos(ctx context.Context, id string) ([]string, error)
	GetPreviewPhoto(ctx context.Context, id string) (string, error)
	PutObject(ctx context.Context, key string, content []byte) (*v4.PresignedHTTPRequest, error)
	DeleteObject(ctx context.Context, key string) (*v4.PresignedHTTPRequest, error)
}
