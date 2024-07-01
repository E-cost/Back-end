package storage

import (
	"Ecost/internal/contact/model"
	"Ecost/internal/utils/ip"
	"Ecost/internal/utils/types"
	"context"
)

type ContactsRepository interface {
	Create(ctx context.Context, contact *model.Contact) error
	GetContactByEmail(ctx context.Context, email string) (string, error)
	CheckEmailExists(ctx context.Context, email string) (string, error)
	GetContactByID(ctx context.Context, id string) (types.ContactInfo, error)
	CreateConfirmationMail(ctx context.Context, verify *model.Confirmation) error
	VerifySecretCode(ctx context.Context, secretCode string) error
	VerifyIsUsed(ctx context.Context, id string) (bool, error)
	UpdateConfirmation(ctx context.Context, id string, ip ip.IPOutput) error
	DeleteExpired(ctx context.Context) (map[string][]string, error)
}
