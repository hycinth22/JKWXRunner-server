package model

import (
	"time"
)

type Ticket struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`

	Contact string `json:"contact"`
	Memo    string `json:"memo"`
	OwnerID uint   `gorm:"INDEX:owner" json:"ownerID"`
}
