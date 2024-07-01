package service

import (
	"Ecost/internal/apperror"
	"Ecost/internal/contact/dto"
	"Ecost/internal/contact/mail"
	"Ecost/internal/contact/model"
	"Ecost/internal/contact/service/helpers"
	"Ecost/internal/contact/storage"
	"Ecost/internal/facades"
	"Ecost/internal/utils"
	"Ecost/internal/utils/ip"
	"Ecost/internal/utils/types"
	"Ecost/pkg/logging"
	"context"
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

type contactService struct {
	logger       *logging.Logger
	cron         *cron.Cron
	mailer       mail.GmailSender
	repository   storage.ContactsRepository
	yandexFacade facades.YandexOrderFacade
}

var (
	instance *contactService
	once     sync.Once
)

func NewContactService(
	logger *logging.Logger,
	repository storage.ContactsRepository,
	yandexFacade facades.YandexOrderFacade,
) ContactService {
	once.Do(func() {
		instance = &contactService{
			logger:       logger,
			cron:         cron.New(),
			repository:   repository,
			yandexFacade: yandexFacade,
		}

		instance.initCronJobs()
	})
	return instance
}

func (s *contactService) SetYandexFacade(facade facades.YandexOrderFacade) {
	s.yandexFacade = facade
}

func (s *contactService) CreateContact(ctx context.Context, dto dto.CreateContactDto, ip ip.IPOutput) (string, string, error) {
	if err := dto.Validate(); err != nil {
		return "", "", fmt.Errorf("validation error: %v", err)
	}

	id, _ := s.repository.CheckEmailExists(ctx, dto.Email)
	if id != "" {
		return "", "", apperror.Conflict("Already exists", "email already exists")
	}

	newContact := &model.Contact{
		Name:         dto.Name,
		Surname:      dto.Surname,
		Email:        dto.Email,
		ContactPhone: dto.ContactPhone,
		Message:      dto.Message,
		IpAddress:    ip.RealIP,
		Port:         ip.Port,
		ProxyChain:   ip.ProxyChain,
	}

	if err := s.repository.Create(ctx, newContact); err != nil {
		return "", "", apperror.InternalServerError("Failed to create contact", err.Error())
	}

	newMail := &model.Confirmation{
		ContactId:  newContact.ID,
		SecretCode: utils.RandomString(32),
		IpAddress:  ip.RealIP,
		Port:       ip.Port,
		ProxyChain: ip.ProxyChain,
	}

	if err := s.repository.CreateConfirmationMail(ctx, newMail); err != nil {
		return "", "", apperror.InternalServerError("Failed to create a mail", err.Error())
	}

	verifyURL := fmt.Sprintf(
		"http://localhost:3000/contacts/confirmation?email_id=%s&secret_code=%s",
		newMail.ID,
		newMail.SecretCode)
	to := []string{dto.Email}
	subject, msg, to, err := helpers.EmailTemplate(dto, verifyURL)
	if err != nil {
		s.logger.Errorf("Failed to generate email template: %v", err)
		return "", "", err
	}

	err = s.mailer.SendMail(subject, string(msg), to, nil, nil, nil)
	if err != nil {
		s.logger.Errorf("Failed to send email: %v", err)
		return "", "", err
	}

	return newMail.ID, newMail.SecretCode, nil
}

func (s *contactService) VerifyContact(ctx context.Context, id, secretCode string, ip ip.IPOutput) error {
	if err := s.repository.VerifySecretCode(ctx, secretCode); err != nil {
		return apperror.BadRequest("Invalid secret code", err.Error())
	}

	if value, err := s.repository.VerifyIsUsed(ctx, secretCode); err != nil {
		return apperror.InternalServerError("Failed to verify the contact", err.Error())
	} else if value {
		return apperror.BadRequest("Invalid or already used secret code", "")
	}

	if err := s.repository.UpdateConfirmation(ctx, id, ip); err != nil {
		return apperror.InternalServerError("Failed to confirm the contact", err.Error())
	}

	return nil
}

func (s *contactService) GetContactByEmail(ctx context.Context, email string) (string, error) {
	id, err := s.repository.GetContactByEmail(ctx, email)
	if err != nil {
		return "", apperror.NotFound("Email is not found", err.Error())
	}

	return id, nil
}

func (s *contactService) GetContactById(ctx context.Context, id string) (types.ContactInfo, error) {
	contactInfo, err := s.repository.GetContactByID(ctx, id)
	if err != nil {
		return types.ContactInfo{}, apperror.NotFound("Contact is not found", err.Error())
	}

	return contactInfo, nil
}

func (s *contactService) DeleteNotVerifiedContact(ctx context.Context) error {
	deletedIds, err := s.repository.DeleteExpired(ctx)
	if err != nil {
		return apperror.InternalServerError("Failed to delete expired contact", err.Error())
	}

	var keys []string
	for _, orders := range deletedIds {
		for _, order := range orders {
			key := fmt.Sprintf("orders/%s.pdf", order)
			keys = append(keys, key)
		}
	}

	if len(keys) > 0 {
		err = s.yandexFacade.DeleteOrder(ctx, keys)
		if err != nil {
			return apperror.InternalServerError("Failed to delete orders from Yandex Object Storage", err.Error())
		}
	}

	return nil
}

func (s *contactService) initCronJobs() {
	_, err := s.cron.AddFunc("@every 8h", func() {
		ctx := context.Background()
		if err := s.DeleteNotVerifiedContact(ctx); err != nil {
			s.logger.Errorf("Failed to delete not verified contact: %v", err)
		}
	})
	if err != nil {
		s.logger.Errorf("Error scheduling contact cron job: %v", err)
	}

	s.cron.Start()
}

func (s *contactService) StopCronJobs() {
	s.cron.Stop()
}
