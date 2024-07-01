package contactsRepository

import (
	"Ecost/internal/contact/model"
	"Ecost/internal/utils/ip"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	confirmation model.Confirmation
)

func (r *repository) CreateConfirmationMail(ctx context.Context, confirmation *model.Confirmation) error {
	q := `
		INSERT INTO confirmations 
			(contact_id, secret_code, is_used, ip_address, port, proxy_chain) 
		VALUES 
			($1, $2, $3, $4, $5, $6) 
		RETURNING id
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if err := r.client.QueryRow(
		ctx,
		q,
		confirmation.ContactId,
		confirmation.SecretCode,
		confirmation.IsUsed,
		confirmation.IpAddress,
		confirmation.Port,
		pq.Array(confirmation.ProxyChain)).Scan(&confirmation.ID); err != nil {
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

func (r *repository) VerifySecretCode(ctx context.Context, secretCode string) error {
	q := `
		SELECT id
		FROM confirmations
		WHERE secret_code = $1
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	err := r.client.QueryRow(ctx, q, secretCode).Scan(&confirmation.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("secret_code has not been found")
		}
		return err
	}

	return nil
}

func (r *repository) VerifyIsUsed(ctx context.Context, secretCode string) (bool, error) {
	q := `
		SELECT is_used
		FROM confirmations
		WHERE secret_code = $1
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var isUsed bool

	err := r.client.QueryRow(ctx, q, secretCode).Scan(&isUsed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no record found for secret code: %s", secretCode)
		}

		return false, err
	}

	return isUsed, nil
}

func (r *repository) UpdateConfirmation(ctx context.Context, id string, ip ip.IPOutput) error {
	q := `
		UPDATE confirmations
        SET is_used = true,
            ip_address = $2,
            port = $3,
            proxy_chain = $4,
            confirmed_at = (now() at time zone 'Europe/Moscow')
		WHERE id = $1
		AND expired_at > now()
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(
		ctx,
		q,
		id,
		ip.RealIP,
		ip.Port,
		ip.ProxyChain)
	if err != nil {
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
