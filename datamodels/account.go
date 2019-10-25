package datamodels

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model

	SchoolID int64  `gorm:"unique_index:schoolID_stuNum;NOT NULL"`
	StuNum   string `gorm:"unique_index:schoolID_stuNum;NOT NULL"`
	Password string `gorm:"NOT NULL"`
	Memo     string `gorm:"NOT NULL;default:''"`

	OwnerID          int          `gorm:"NOT NULL;default:0"`
	DeviceID         uint         `gorm:"NOT NULL"`
	Status           string       `gorm:"NOT NULL;default:'normal'"`
	RunDistance      float64      `gorm:"NOT NULL"`
	StartDistance    float64      `gorm:"NOT NULL"`
	FinishDistance   float64      `gorm:"NOT NULL"`
	CheckCheatMarked sql.NullBool `gorm:"NOT NULL;default:1"`

	LastResult sql.NullString `gorm:"default:''"`
	LastTime   sql.NullTime   `gorm:"default:0"`
}
