package model

import (
	"time"
)

type Order struct {
	ID             string    `json:"id"`
	ContactId      string    `json:"contact_id"`
	ShortId        string    `json:"short_id"`
	Email          string    `json:"email"`
	Items          []Item    `json:"items"`
	TotalAmountByn float64   `json:"total_amount_byn"`
	IpAddress      string    `json:"ip_address"`
	Port           string    `json:"port"`
	ProxyChain     []string  `json:"proxy_chain"`
	CreatedAt      time.Time `json:"created_at"`
}

type Item struct {
	ID        string  `json:"id"`
	Article   string  `json:"article"`
	Category  string  `json:"category"`
	Product   string  `json:"product"`
	Name      string  `json:"name"`
	NetWeight string  `json:"net_weight"`
	Quantity  int     `json:"quantity"`
	PriceBYN  float32 `json:"price_byn"`
}
