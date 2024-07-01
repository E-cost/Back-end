package contactsRepository

import (
	"Ecost/internal/contact/model"
	"Ecost/internal/utils/types"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (r *repository) Create(ctx context.Context, contact *model.Contact) error {
	q := `
			INSERT INTO contacts 
			    (name, surname, email, contact_phone, message, ip_address, port, proxy_chain) 
			VALUES 
			    ($1, $2, $3, $4, $5, $6, $7, $8) 
			RETURNING id
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if err := r.client.QueryRow(
		ctx,
		q,
		contact.Name,
		contact.Surname,
		contact.Email,
		contact.ContactPhone,
		contact.Message,
		contact.IpAddress,
		contact.Port,
		contact.ProxyChain).Scan(&contact.ID); err != nil {
		if errors.As(err, &pgErr) {
			errors.As(err, &pgErr)
			newErr := fmt.Errorf(fmt.Sprintf(
				"SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message,
				pgErr.Detail,
				pgErr.Where,
				pgErr.Code,
				pgErr.SQLState()))

			r.logger.Error(newErr)

			return newErr
		}

		return err
	}

	return nil
}

func (r *repository) GetContactByEmail(ctx context.Context, email string) (string, error) {
	q := `
		SELECT id
		FROM contacts
		WHERE email = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var contactID string

	err := r.client.QueryRow(ctx, q, email).Scan(&contactID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no record found for contact email: %s", email)
		}

		return "", err
	}

	return contactID, nil
}

func (r *repository) CheckEmailExists(ctx context.Context, email string) (string, error) {
	q := `
		SELECT id
		FROM contacts
		WHERE email = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var contactID string

	err := r.client.QueryRow(ctx, q, email).Scan(&contactID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Trace("No contact found with the given email")
			return "", nil
		}
		r.logger.Errorf("Failed to execute query: %v", err)
		return "", err
	}

	return contactID, nil
}

func (r *repository) GetContactByID(ctx context.Context, id string) (types.ContactInfo, error) {
	q := `
		SELECT name, surname, email, contact_phone
		FROM contacts
		WHERE id = $1
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var contact types.ContactInfo

	err := r.client.QueryRow(ctx, q, id).Scan(&contact.Name, &contact.Surname, &contact.Email, &contact.ContactPhone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.ContactInfo{}, fmt.Errorf("no record found for contact: %s", id)
		}

		return types.ContactInfo{}, err
	}

	return contact, nil
}

func (r *repository) DeleteExpired(ctx context.Context) (map[string][]string, error) {
	r.logger.Info("Cleaning up expired contacts")

	selectQuery := `
		SELECT 
		    contacts.id, 
		    confirmations.id, 
			orders.short_id,
		    confirmations.is_used,
		    confirmations.expired_at,
		    now() AT TIME ZONE 'Europe/Moscow' as "current_time"
		FROM contacts
		INNER JOIN confirmations
			ON contacts.id = confirmations.contact_id
		INNER JOIN orders
			ON contacts.id = orders.contact_id
		WHERE confirmations.expired_at <= (now() AT TIME ZONE 'Europe/Moscow') AND confirmations.is_used = false;
	`

	rows, err := r.client.Query(ctx, selectQuery)
	if err != nil {
		r.logger.Errorf("Failed to delete expired contacts: %v", err.Error())
		return nil, err
	}
	defer rows.Close()

	contactOrdersMap := make(map[string][]string)

	for rows.Next() {
		var (
			contactId      string
			confirmationId string
			orderShortId   string
			expiredAt      time.Time
			isUsed         bool
			currentTime    time.Time
		)
		if err := rows.Scan(&contactId, &confirmationId, &orderShortId, &isUsed, &expiredAt, &currentTime); err != nil {
			r.logger.Errorf("Failed to scan deleted contact: %v", err.Error())
			return nil, err
		}
		contactOrdersMap[contactId] = append(contactOrdersMap[contactId], orderShortId)
		r.logger.Infof(
			"Deleted contact: id=%s, confirmation_id=%s, expired_at=%s, is_used=%v, current_time=%v",
			contactId, confirmationId, expiredAt, isUsed, currentTime)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error iterating over deleted contacts: %v", err)
		return nil, err
	}

	if len(contactOrdersMap) == 0 {
		r.logger.Info("No expired contacts found")
		return nil, nil
	}

	deleteQuery := `
		DELETE FROM contacts
		WHERE id IN (
			SELECT contacts.id
			FROM confirmations
			INNER JOIN contacts
				ON confirmations.contact_id = contacts.id
			WHERE confirmations.expired_at <= (now() AT TIME ZONE 'Europe/Moscow') AND confirmations.is_used = false
		);
	`

	_, err = r.client.Exec(ctx, deleteQuery)
	if err != nil {
		r.logger.Errorf("Failed to delete expired contacts: %v", err.Error())
		return nil, err
	}

	r.logger.Info("Expired contacts deleted successfully")
	return contactOrdersMap, nil
}
