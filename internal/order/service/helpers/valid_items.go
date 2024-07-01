package helpers

import (
	"Ecost/internal/facades"
	"Ecost/internal/order/model"
)

func GetValidItems(results []facades.ItemResult) []model.Item {
	validItems := make([]model.Item, 0)
	for _, result := range results {
		if result.Valid {
			validItems = append(validItems, result.Item)
		}
	}

	return validItems
}
