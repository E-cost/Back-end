package service

import (
	"Ecost/internal/apperror"
	"Ecost/internal/config"
	"Ecost/pkg/logging"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type yandexService struct {
	logger *logging.Logger
	client *s3.PresignClient
	config *config.Config
}

func NewYandexService(
	logger *logging.Logger,
	client *s3.PresignClient,
	config *config.Config,
) YandexService {
	return &yandexService{
		logger: logger,
		client: client,
		config: config,
	}
}

func (s *yandexService) DownloadPriceList(ctx context.Context) ([]byte, string, error) {
	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	key := "price-list/price-list.xlsx"

	req, err := s.client.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Yandex.Bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(1 * time.Hour)
	})
	if err != nil {
		return nil, "", apperror.InternalServerError("failed to generate presigned URL", err.Error())
	}

	resp, err := http.Get(req.URL)
	if err != nil {
		return nil, "", apperror.InternalServerError("failed to download price list using presigned URL", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", apperror.InternalServerError("failed to download price list using presigned URL", fmt.Sprintf("status code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", apperror.InternalServerError("failed to read price list from response body", err.Error())
	}

	return body, contentType, nil
}

func (s *yandexService) GetPreviewPhoto(ctx context.Context, id string) (string, error) {
	key := fmt.Sprintf("photos/%s/preview.png", id)

	request, err := s.client.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.config.Yandex.Bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(7 * 24 * time.Hour)
	})
	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func (s *yandexService) GetItemPhotos(ctx context.Context, id string) ([]string, error) {
	urls := make([]string, 0, 3)

	for i := 1; i <= 3; i++ {
		key := fmt.Sprintf("photos/%s/photo_%d.png", id, i)

		request, err := s.client.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(s.config.Yandex.Bucket),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(7 * 24 * time.Hour)
		})
		if err != nil {
			return nil, err
		}

		urls = append(urls, request.URL)
	}

	return urls, nil
}

func (s *yandexService) PutObject(ctx context.Context, key string, content []byte) (*v4.PresignedHTTPRequest, error) {
	req, err := s.client.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.config.Yandex.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(content),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(30 * time.Minute)
	})
	if err != nil {
		s.logger.Errorf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n", s.config.Yandex.Bucket, key, err)
		return nil, apperror.InternalServerError("Couldn't get a presigned request to put an object. Here's why: %v\n", err.Error())
	}

	return req, nil
}

func (s *yandexService) DeleteObject(ctx context.Context, key string) (*v4.PresignedHTTPRequest, error) {
	req, err := s.client.PresignDeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.Yandex.Bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(30 * time.Minute)
	})
	if err != nil {
		s.logger.Errorf("Couldn't get a presigned request to delete %v:%v. Here's why: %v\n", s.config.Yandex.Bucket, key, err)
		return nil, apperror.InternalServerError("Couldn't get a presigned request to delete an object. Here's why: %v\n", err.Error())
	}

	return req, nil
}
