package service

import (
	"errors"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
)

var (
	ErrNoDevice = errors.New("没有找到该用户的设备")
)

type IDeviceService interface {
	ICommonService
	GetDevice(deviceID uint) (datamodels.Device, error)
	SaveDevice(device *datamodels.Device) error
}

type deviceService struct {
	ICommonService
	db database.TX
}

func (d *deviceService) GetDevice(deviceID uint) (datamodels.Device, error) {
	device := datamodels.Device{}
	device.ID = deviceID
	if err := d.db.First(&device).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			return device, ErrNoDevice
		}
		return device, WrapAsInternalError(err)
	}
	return device, nil
}

func (d *deviceService) SaveDevice(device *datamodels.Device) error {
	err := d.db.Save(device).Error
	if err != nil {
		return WrapAsInternalError(err)
	}
	return nil
}

func NewDeviceService() IDeviceService {
	return NewDeviceServiceUpon(NewCommonService())
}

func NewDeviceServiceOn(db *database.DB) IDeviceService {
	return NewDeviceServiceUpon(NewCommonServiceOn(db))
}

func NewDeviceServiceUpon(commonService ICommonService) IDeviceService {
	return &deviceService{ICommonService: commonService, db: commonService.GetDB()}
}
