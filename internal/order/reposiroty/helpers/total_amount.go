package helpers

import "Ecost/internal/order/model"

func CalculateTotalAmount(items []model.Item) float64 {
	total := 0.0
	for _, item := range items {
		total += float64(item.Quantity) * float64(item.PriceBYN)
	}
	return total
}
