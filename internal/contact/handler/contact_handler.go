package contactsHandler

import (
	"Ecost/internal/api/handlers"
	"Ecost/internal/api/middleware/dns"
	"Ecost/internal/api/middleware/errors"
	"Ecost/internal/api/middleware/limit"
	"Ecost/internal/api/middleware/verify"
	"Ecost/internal/apperror"
	"Ecost/internal/contact/dto"
	"Ecost/internal/contact/service"
	"Ecost/internal/contact/storage"
	"Ecost/internal/facades"
	"Ecost/internal/utils/ip"
	"Ecost/pkg/logging"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	GETConfirmationURL   = "/contacts/confirmation"
	POSTContactCreateURL = "/contacts/send"
)

var ()

type handler struct {
	logger       *logging.Logger
	ipMiddleware *dns.Middleware
	service      service.ContactService
	repository   storage.ContactsRepository
}

func NewContactHandler(
	logger *logging.Logger,
	ipMiddleware *dns.Middleware,
	repository storage.ContactsRepository,
	yandexFacade facades.YandexOrderFacade,
) handlers.ContactHandler {

	contactService := service.NewContactService(logger, repository, yandexFacade)
	return &handler{
		logger:       logger,
		ipMiddleware: ipMiddleware,
		service:      contactService,
		repository:   repository,
	}
}

func (h *handler) SaveContact(router *httprouter.Router) {
	router.HandlerFunc(
		http.MethodGet,
		GETConfirmationURL,
		limit.WriteMiddleware(verify.Middleware(errors.Middleware(h.VerifyContact))))
	router.HandlerFunc(
		http.MethodPost,
		POSTContactCreateURL,
		limit.WriteMiddleware(h.ipMiddleware.IpMiddleware(errors.Middleware(h.CreateContact))))
}

func (h *handler) CreateContact(w http.ResponseWriter, r *http.Request) error {
	var createContactDto dto.CreateContactDto

	IPOutput, err := ip.ReadUserIP(r)
	if err != nil {
		return apperror.InternalServerError("Failed to get ip from request due an error:", err.Error())
	}

	if err := json.NewDecoder(r.Body).Decode(&createContactDto); err != nil {
		return apperror.BadRequest("Invalid request body", err.Error())
	}

	emailID, secretCode, err := h.service.CreateContact(r.Context(), createContactDto, *IPOutput)
	if err != nil {
		return apperror.BadRequest("Failed to create contact", err.Error())
	}

	response := struct {
		EmailID    string `json:"email_id"`
		SecretCode string `json:"secret_code"`
	}{
		EmailID:    emailID,
		SecretCode: secretCode,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return apperror.InternalServerError("Failed to encode response", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *handler) VerifyContact(w http.ResponseWriter, r *http.Request) error {
	emailIDStr := r.URL.Query().Get("email_id")
	secretCode := r.URL.Query().Get("secret_code")

	IPOutput, err := ip.ReadUserIP(r)
	if err != nil {
		return apperror.InternalServerError("Failed to get ip from request due an error:", err.Error())
	}

	verifyErr := h.service.VerifyContact(r.Context(), emailIDStr, secretCode, *IPOutput)
	if verifyErr != nil {
		return apperror.InternalServerError("Can not verify the contact", verifyErr.Error())
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
