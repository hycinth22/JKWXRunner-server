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

	LimitRandMin  float64 `gorm:"NOT NULL"`
	LimitRandMax  float64 `gorm:"NOT NULL"`
	LimitDailyMin float64 `gorm:"NOT NULL"`
	LimitDailyMax float64 `gorm:"NOT NULL"`

	DeviceID   uint      `gorm:"NOT NULL"`
	Status     uint      `gorm:"NOT NULL"`
	LastResult uint      `gorm:"NOT NULL"`
	LastTime   time.Time `gorm:"NOT NULL"`
}
