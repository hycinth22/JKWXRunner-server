package model

import (
	"time"
)

type Token struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	RemoteUserID   int64     `gorm:"primary_key;NOT NULL"`
	TokenID        string    `gorm:"NOT NULL"`
	ExpirationTime time.Time `gorm:"NOT NULL"`
}
