// 提供对Device的管理
package deviceSrv

import (
	"errors"
	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/database/model"
	"github.com/inkedawn/JKWXRunner-server/service"
	"github.com/inkedawn/go-sunshinemotion"
)

type Device = model.Device

var (
	ErrNoDevice = errors.New("没有找到该用户的设备")
)

func GetDevice(db *database.DB, deviceID uint) (Device, error) {
	device := model.Device{}
	device.ID = deviceID
	if err := db.First(&device).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			return device, ErrNoDevice
		}
		return device, service.WrapAsInternalError(err)
	}
	return device, nil
}

func SaveDevice(db *database.DB, device *Device) error {
	err := db.Save(device).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	return nil
}

func FromSSMTDevice(device ssmt.Device) Device {
	return Device{
		DeviceName: device.DeviceName,
		ModelType:  device.ModelType,
		Screen:     device.Screen,
		IMEI:       device.IMEI,
		IMSI:       device.IMSI,
		UserAgent:  device.UserAgent,
	}
}
func ToSSMTDevice(device Device) ssmt.Device {
	return ssmt.Device{
		DeviceName: device.DeviceName,
		ModelType:  device.ModelType,
		Screen:     device.Screen,
		IMEI:       device.IMEI,
		IMSI:       device.IMSI,
		UserAgent:  device.UserAgent,
	}
}
