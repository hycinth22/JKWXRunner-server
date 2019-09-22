package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model

	SchoolID int64  `gorm:"unique_index:schoolID_stuNum;NOT NULL"`
	StuNum   string `gorm:"unique_index:schoolID_stuNum;NOT NULL"`
	Password string `gorm:"NOT NULL"`
	Memo     string `gorm:"NOT NULL;default:''"`

	OwnerID          int     `gorm:"NOT NULL"`
	DeviceID         uint    `gorm:"NOT NULL"`
	Status           string  `gorm:"NOT NULL;default:'normal'"`
	RunDistance      float64 `gorm:"NOT NULL"`
	StartDistance    float64 `gorm:"NOT NULL"`
	FinishDistance   float64 `gorm:"NOT NULL"`
	CheckCheatMarked bool    `gorm:"NOT NULL;default:1"`

	LastResult string    `gorm:"default:''"`
	LastTime   time.Time `gorm:"default:0"`
}
