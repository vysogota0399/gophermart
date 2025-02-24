package models

import "time"

type Session struct {
	ID        uint64    `json:"id"`
	Sub       string    `json:"sub"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}
