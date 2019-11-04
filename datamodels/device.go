package datamodels

import (
	ssmt "github.com/inkedawn/go-sunshinemotion/v3"
	"github.com/jinzhu/gorm"
)

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

func DeviceFromSSMTDevice(device ssmt.Device) Device {
	return Device{
		DeviceName: device.DeviceName,
		ModelType:  device.ModelType,
		Screen:     device.Screen,
		IMEI:       device.IMEI,
		IMSI:       device.IMSI,
		UserAgent:  device.UserAgent,
	}
}
func DeviceToSSMTDevice(device Device) ssmt.Device {
	return ssmt.Device{
		DeviceName: device.DeviceName,
		ModelType:  device.ModelType,
		Screen:     device.Screen,
		IMEI:       device.IMEI,
		IMSI:       device.IMSI,
		UserAgent:  device.UserAgent,
	}
}
