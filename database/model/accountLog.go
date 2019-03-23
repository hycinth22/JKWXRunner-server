package model

import (
	"time"
)

type AccountLog struct {
	UID     uint      `gorm:"index"`
	Time    time.Time `gorm:"NOT NULL"`
	Type    uint      `gorm:"NOT NULL"`
	Content string    `gorm:"NOT NULL"`
}
