package helpers

import (
	"Ecost/internal/item/model"
	"Ecost/internal/utils/types"
)

func ConvertToSeafoodType(item model.Seafood) (*types.GetSeafoodType, error) {
	itemType := &types.GetSeafoodType{
		ID:             item.ID,
		Article:        item.Article,
		Category:       item.Category,
		Product:        item.Product,
		Name:           item.Name,
		CountryId:      item.CountryId,
		NetWeight:      item.NetWeight,
		Composition:    item.Composition,
		FoodValue:      item.FoodValue,
		Supplements:    item.Supplements,
		Vitamins:       item.Vitamins,
		EnergyValue:    item.EnergyValue,
		Description:    item.Description,
		Recommendation: item.Recommendation,
		ShelfLife:      item.ShelfLife,
		ExpirationDate: item.ExpirationDate,
		PriceBYN:       item.PriceBYN,
	}

	return itemType, nil
}

func ConvertToSeafoodSliceType(items []model.Seafood, urls map[string]string) ([]*types.GetSeafoodType, error) {
	itemsType := make([]*types.GetSeafoodType, len(items))
	for i, item := range items {
		itemType, err := ConvertToSeafoodType(item)
		if err != nil {
			return nil, err
		}
		if url, exists := urls[item.ID]; exists {
			itemType.PreviewUrl = url
		}
		itemsType[i] = itemType
	}

	return itemsType, nil
}
