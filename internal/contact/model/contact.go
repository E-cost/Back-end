package model

import (
	"time"

	"github.com/volatiletech/null/v8"
)

type Contact struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Email        string    `json:"email"`
	ContactPhone string    `json:"contact_phone"`
	Message      string    `json:"message"`
	IsValid      bool      `json:"is_valid"`
	IpAddress    string    `json:"ip_address"`
	Port         string    `json:"port"`
	ProxyChain   []string  `json:"proxy_chain"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    null.Time `json:"updated_at"`
}
