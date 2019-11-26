package datamodels

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model

	SchoolID int64  `gorm:"index:personalID;NOT NULL"`
	StuNum   string `gorm:"index:personalID;NOT NULL"`
	Password string `gorm:"NOT NULL"`
	Memo     string `gorm:"NOT NULL;default:''"`

	OwnerID          int          `gorm:"NOT NULL;default:0"`
	DeviceID         uint         `gorm:"NOT NULL"`
	Status           string       `gorm:"NOT NULL;default:'normal'"`
	RunDistance      float64      `gorm:"NOT NULL"`
	StartDistance    float64      `gorm:"NOT NULL"`
	FinishDistance   float64      `gorm:"NOT NULL"`
	CheckCheatMarked sql.NullBool `gorm:"NOT NULL;default:1"`

	LastResult sql.NullString
	LastTime   sql.NullTime
}

// Model that all filed is empty
var AccountModel = new(Account)
