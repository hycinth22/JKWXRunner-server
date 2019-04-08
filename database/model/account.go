package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Account struct {
	gorm.Model

	SchoolID int64  `gorm:"unique_index:schoolID_stuNum;NOT NULL"`
	StuNum   string `gorm:"unique_index:schoolID_stuNum;NOT NULL"`
	Password string `gorm:"NOT NULL"`

	RunDistance      float64 `gorm:"NOT NULL"`
	CheckCheatMarked bool    `gorm:"NOT NULL;default:true"`

	DeviceID   uint      `gorm:"NOT NULL"`
	Status     string    `gorm:"NOT NULL"`
	LastResult string    `gorm:"default:''"`
	LastTime   time.Time `gorm:"default:0"`
}
