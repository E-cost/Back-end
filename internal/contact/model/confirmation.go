package model

import (
	"time"
)

type Confirmation struct {
	ID         string    `json:"id"`
	ContactId  string    `json:"contact_id"`
	SecretCode string    `json:"secret_code"`
	IsUsed     bool      `json:"is_used"`
	IpAddress  string    `json:"ip_address"`
	Port       string    `json:"port"`
	ProxyChain []string  `json:"proxy_chain"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}
