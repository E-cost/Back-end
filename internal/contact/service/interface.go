package service

import (
	"Ecost/internal/contact/dto"
	"Ecost/internal/facades"
	"Ecost/internal/utils/ip"
	"Ecost/internal/utils/types"
	"context"
)

type ContactService interface {
	SetYandexFacade(facade facades.YandexOrderFacade)
	CreateContact(ctx context.Context, dto dto.CreateContactDto, ip ip.IPOutput) (string, string, error)
	VerifyContact(ctx context.Context, id, secretCode string, ip ip.IPOutput) error
	GetContactByEmail(ctx context.Context, email string) (string, error)
	GetContactById(ctx context.Context, id string) (types.ContactInfo, error)
	DeleteNotVerifiedContact(ctx context.Context) error
	initCronJobs()
	StopCronJobs()
}
