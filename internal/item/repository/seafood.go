package itemsRepository

import (
	"Ecost/internal/item/model"
	"Ecost/internal/utils/sort"
	"context"
	"fmt"
)

func (r *repository) FindOneSeafood(ctx context.Context, id string) (model.Seafood, error) {
	q := `
		SELECT *
		FROM seafood
		WHERE id = $1
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var seafood model.Seafood

	err := r.client.QueryRow(ctx, q, id).Scan(
		&seafood.ID,
		&seafood.Article,
		&seafood.Category,
		&seafood.Product,
		&seafood.Name,
		&seafood.CountryId,
		&seafood.NetWeight,
		&seafood.Composition,
		&seafood.Supplements,
		&seafood.FoodValue,
		&seafood.Vitamins,
		&seafood.EnergyValue,
		&seafood.Description,
		&seafood.Recommendation,
		&seafood.ShelfLife,
		&seafood.ExpirationDate,
		&seafood.PriceBYN,
		&seafood.PriceUSD,
		&seafood.CreatedAt,
		&seafood.UpdatedAt)
	if err != nil {
		return model.Seafood{}, err
	}

	return seafood, nil
}

func (r *repository) FindAllSeafood(ctx context.Context, sortOptions sort.SortOptions) ([]model.Seafood, error) {
	orderBy := sortOptions.GetOrderBy()

	q := fmt.Sprintf(`
		SELECT *
		FROM seafood
		ORDER BY %s;
	`, orderBy)

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var seafoodItems []model.Seafood
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var seafood model.Seafood

		err = rows.Scan(
			&seafood.ID,
			&seafood.Article,
			&seafood.Category,
			&seafood.Product,
			&seafood.Name,
			&seafood.CountryId,
			&seafood.NetWeight,
			&seafood.Composition,
			&seafood.Supplements,
			&seafood.FoodValue,
			&seafood.Vitamins,
			&seafood.EnergyValue,
			&seafood.Description,
			&seafood.Recommendation,
			&seafood.ShelfLife,
			&seafood.ExpirationDate,
			&seafood.PriceBYN,
			&seafood.PriceUSD,
			&seafood.CreatedAt,
			&seafood.UpdatedAt)
		if err != nil {
			return nil, err
		}

		seafoodItems = append(seafoodItems, seafood)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return seafoodItems, nil
}

func (r *repository) GetAllSeafood(ctx context.Context) ([]model.Seafood, error) {
	q := `
		SELECT *
		FROM seafood;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var seafoodItems []model.Seafood
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var seafood model.Seafood

		err = rows.Scan(
			&seafood.ID,
			&seafood.Article,
			&seafood.Category,
			&seafood.Product,
			&seafood.Name,
			&seafood.CountryId,
			&seafood.NetWeight,
			&seafood.Composition,
			&seafood.Supplements,
			&seafood.FoodValue,
			&seafood.Vitamins,
			&seafood.EnergyValue,
			&seafood.Description,
			&seafood.Recommendation,
			&seafood.ShelfLife,
			&seafood.ExpirationDate,
			&seafood.PriceBYN,
			&seafood.PriceUSD,
			&seafood.CreatedAt,
			&seafood.UpdatedAt)
		if err != nil {
			return nil, err
		}

		seafoodItems = append(seafoodItems, seafood)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return seafoodItems, nil
}

func (r *repository) GetAllIdsSeafood(ctx context.Context) ([]string, error) {
	q := `
		SELECT id
		FROM seafood;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	var seafoodIds []string

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var seafoodID string
		err := rows.Scan(&seafoodID)
		if err != nil {
			return nil, err
		}
		seafoodIds = append(seafoodIds, seafoodID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return seafoodIds, nil
}
