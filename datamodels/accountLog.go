package datamodels

import (
	"time"
)

type AccountLog struct {
	UID     uint      `gorm:"index"`
	Time    time.Time `gorm:"NOT NULL"`
	Type    string    `gorm:"NOT NULL"`
	Content string    `gorm:"NOT NULL"`
}
