package model

import (
	"time"

	"github.com/volatiletech/null/v8"
)

type Seafood struct {
	ID             string    `json:"id"`
	Article        string    `json:"article"`
	Category       string    `json:"category"`
	Product        string    `json:"product"`
	Name           string    `json:"name"`
	CountryId      string    `json:"country_id"`
	NetWeight      string    `json:"net_weight"`
	Composition    *string   `json:"composition"`
	FoodValue      *string   `json:"food_value"`
	Supplements    *string   `json:"supplements"`
	Vitamins       *string   `json:"vitamins"`
	EnergyValue    *string   `json:"energy_value"`
	Description    *string   `json:"description"`
	Recommendation *string   `json:"recommendation"`
	ShelfLife      *string   `json:"shelf_life"`
	ExpirationDate *string   `json:"expiration_date"`
	PriceBYN       *float32  `json:"price_byn"`
	PriceUSD       *float32  `json:"price_usd"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      null.Time `json:"updated_at"`
}
