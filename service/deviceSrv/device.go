// 提供对Device的管理
package deviceSrv

import (
	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
)

// DEPRECATED, use datamodels.Device
type Device = datamodels.Device

var (
	// DEPRECATED, use service.ErrNoDevice
	ErrNoDevice = service.ErrNoDevice
)

// DEPRECATED, use service.IDeviceService instead
func GetDevice(db *database.DB, deviceID uint) (Device, error) {
	return service.NewDeviceServiceOn(db).GetDevice(deviceID)
}

// DEPRECATED, use service.IDeviceService instead
func SaveDevice(db *database.DB, device *Device) error {
	return service.NewDeviceServiceOn(db).SaveDevice(device)
}

// DEPRECATED, use datamodels.DeviceFromSSMTDevice instead
func FromSSMTDevice(device ssmt.Device) Device {
	return datamodels.DeviceFromSSMTDevice(device)
}

// DEPRECATEDuse datamodels.DeviceToSSMTDevice instead
func ToSSMTDevice(device Device) ssmt.Device {
	return datamodels.DeviceToSSMTDevice(device)
}
