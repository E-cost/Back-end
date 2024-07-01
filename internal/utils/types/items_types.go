package types

type AllItemsResponse struct {
	Seafood []*GetSeafoodType `json:"seafood"`
}

type GetSeafoodType struct {
	ID             string   `json:"id"`
	Article        string   `json:"article"`
	Category       string   `json:"category"`
	Product        string   `json:"product"`
	Name           string   `json:"name"`
	CountryId      string   `json:"country_id"`
	NetWeight      string   `json:"net_weight"`
	Composition    *string  `json:"composition"`
	FoodValue      *string  `json:"food_value"`
	Supplements    *string  `json:"supplements"`
	Vitamins       *string  `json:"vitamins"`
	EnergyValue    *string  `json:"energy_value"`
	Description    *string  `json:"description"`
	Recommendation *string  `json:"recommendation"`
	ShelfLife      *string  `json:"shelf_life"`
	ExpirationDate *string  `json:"expiration_date"`
	PriceBYN       *float32 `json:"price_byn"`
	PreviewUrl     string   `json:"preview_url"`
}
