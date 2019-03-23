package model

import "github.com/jinzhu/gorm"

//noinspection SpellCheckingInspection
type Device struct {
	gorm.Model

	DeviceName string `gorm:"NOT NULL"`
	ModelType  string `gorm:"NOT NULL"`
	Screen     string `gorm:"NOT NULL"`
	IMEI       string `gorm:"NOT NULL"`
	IMSI       string `gorm:"NOT NULL"`
	UserAgent  string `gorm:"NOT NULL"`
}
