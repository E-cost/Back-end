package ordersRepository

import (
	"Ecost/internal/order/model"
	"Ecost/internal/order/reposiroty/helpers"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func (r *repository) Create(ctx context.Context, order *model.Order) (string, string, float64, error) {
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return "", "", 0.00, fmt.Errorf("failed to serialize items: %v", err)
	}

	totalAmountByn := helpers.CalculateTotalAmount(order.Items)
	order.TotalAmountByn = totalAmountByn

	q := `
        INSERT INTO orders
            (contact_id, items, total_amount_byn, ip_address, port, proxy_chain) 
        VALUES 
            ($1, $2, $3, $4, $5, $6)
        RETURNING id, short_id;
    `

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var orderID string
	var shortID string

	if err := r.client.QueryRow(
		ctx,
		q,
		order.ContactId,
		itemsJSON,
		totalAmountByn,
		order.IpAddress,
		order.Port,
		order.ProxyChain,
	).Scan(&orderID, &shortID); err != nil {
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

			return "", "", 0.00, newErr
		}

		return "", "", 0.00, err
	}

	return orderID, shortID, totalAmountByn, nil
}

func (r *repository) DeleteExpired(ctx context.Context) ([]string, error) {
	r.logger.Info("Cleaning up expired orders")

	selectQuery := `
        SELECT id, short_id, expired_at, is_paid, now() AT TIME ZONE 'Europe/Moscow' as "current_time"
        FROM orders
        WHERE expired_at <= (now() AT TIME ZONE 'Europe/Moscow') AND is_paid = false;
	`

	rows, err := r.client.Query(ctx, selectQuery)
	if err != nil {
		r.logger.Errorf("Failed to delete expired orders: %v", err)
		return nil, err
	}
	defer rows.Close()

	var idsToDelete []string

	for rows.Next() {
		var (
			id          string
			shortId     string
			expiredAt   time.Time
			isPaid      bool
			currentTime time.Time
		)

		if err := rows.Scan(&id, &shortId, &expiredAt, &isPaid, &currentTime); err != nil {
			r.logger.Errorf("Failed to scan deleted order: %v", err)
			return nil, err
		}
		idsToDelete = append(idsToDelete, shortId)
		r.logger.Infof(
			"Deleted order: id=%s, short_id=%s, expired_at=%v, is_paid=%v, current_time=%v",
			id, shortId, expiredAt, isPaid, currentTime)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error iterating over deleted orders: %v", err)
		return nil, err
	}

	if len(idsToDelete) == 0 {
		r.logger.Info("No expired orders found")
		return nil, nil
	}

	deleteQuery := `
        DELETE FROM orders
        WHERE short_id = ANY($1);
    `

	_, err = r.client.Exec(ctx, deleteQuery, idsToDelete)
	if err != nil {
		r.logger.Errorf("Failed to delete expired orders: %v", err)
		return nil, err
	}

	r.logger.Info("Expired orders deleted successfully")
	return idsToDelete, nil
}
