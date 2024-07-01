package service

import (
	"Ecost/internal/apperror"
	dto2 "Ecost/internal/contact/dto"
	service2 "Ecost/internal/contact/service"
	"Ecost/internal/facades"
	"Ecost/internal/order/dto"
	"Ecost/internal/order/service"
	"Ecost/internal/utils/ip"
	"Ecost/internal/utils/types"
	"context"
	"time"
)

type orderContactFacade struct {
	orderService   service.OrderService
	contactService service2.ContactService
}

func NewOrderContactFacade(orderService service.OrderService, contactService service2.ContactService) facades.ContactOrderFacade {
	return &orderContactFacade{
		orderService:   orderService,
		contactService: contactService,
	}
}

func (f *orderContactFacade) CreateContactOrder(ctx context.Context, dto dto.CreateOrderDto, ip ip.IPOutput) (types.ContactConfirmationInfo, error) {
	id, err := f.contactService.GetContactByEmail(ctx, dto.Email)
	output := types.ContactConfirmationInfo{}
	if err != nil {
		if apperror.IsNotFound(err) {
			emailId, secretCode, err := f.contactService.CreateContact(ctx, dto2.CreateContactDto{
				Name:         dto.Name,
				Surname:      dto.Surname,
				Email:        dto.Email,
				ContactPhone: dto.ContactPhone,
			}, ip)
			if err != nil {
				return types.ContactConfirmationInfo{}, apperror.InternalServerError("Failed to create contact", err.Error())
			}

			retries := 3
			for i := 0; i < retries; i++ {
				id, err := f.contactService.GetContactByEmail(ctx, dto.Email)
				if err == nil {
					output = types.ContactConfirmationInfo{
						ContactId:  id,
						EmailId:    emailId,
						SecretCode: secretCode,
					}
					return output, nil
				}
				if !apperror.IsNotFound(err) {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			return output, nil
		}

		return types.ContactConfirmationInfo{}, apperror.InternalServerError("Failed to get contact", err.Error())
	}

	output = types.ContactConfirmationInfo{
		ContactId: id,
	}
	return output, nil
}

func (f *orderContactFacade) GetContactOrder(ctx context.Context, id string) (types.ContactInfo, error) {
	contactInfo, err := f.contactService.GetContactById(ctx, id)
	if err != nil {
		return types.ContactInfo{}, apperror.InternalServerError("Failed to get contact", err.Error())
	}

	return contactInfo, nil
}
