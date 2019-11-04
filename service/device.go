package service

import (
	"errors"
	"sync"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
)

var (
	ErrNoDevice = errors.New("没有找到该用户的设备")
)

type IDeviceService interface {
	GetDevice(deviceID uint) (datamodels.Device, error)
	SaveDevice(device *datamodels.Device) error
}

type deviceService struct {
	db *database.DB
	sync.Locker
}

func (d *deviceService) GetDevice(deviceID uint) (datamodels.Device, error) {
	d.Lock()
	defer d.Unlock()
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
	d.Lock()
	defer d.Unlock()
	err := d.db.Save(device).Error
	if err != nil {
		return WrapAsInternalError(err)
	}
	return nil
}

func NewDeviceService() IDeviceService {
	return NewDeviceServiceOn(database.GetDB())
}

func NewDeviceServiceOn(db *database.DB) IDeviceService {
	return &deviceService{db: db, Locker: &sync.Mutex{}}
}
