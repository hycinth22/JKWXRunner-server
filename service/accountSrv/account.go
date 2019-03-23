// 提供对Account的管理
package accountSrv

import (
	"errors"
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/database/model"
	"github.com/inkedawn/JKWXFucker-server/service"

	"time"
)

type Account = model.Account

var (
	ErrNoAccount = errors.New("没有找到帐号")
)

type Status = uint

const (
	StatusNormal Status = iota
	StatusRunning
	StatusFinished
	StatusSuspend
	StatusTerminated
)

type RunResult = uint

const (
	RunSuccess RunResult = iota
	RunErrorOccurred
)

func ListAccounts(db *database.DB, offset, num uint) ([]Account, error) {
	var accounts []Account
	if err := db.Offset(offset).Limit(num).Find(&accounts).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			// 返回空集
			return accounts, nil
		}
		return accounts, service.WrapAsInternalError(err)
	}
	return accounts, nil
}

func SaveAccount(db *database.DB, acc *Account) error {
	err := db.Save(&acc).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	return nil
}

func ListAndSetRunStatusForAllAccountsWaitRun(db *database.DB) (accounts []Account, err error) {
	now := time.Now()
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	var idGroup []uint
	if err := tx.Model(&Account{}).Where("status = ? AND last_time < ?", StatusNormal, todayZero).Pluck("id", &idGroup).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			// 返回空集
			return []Account{}, nil
		}
		return accounts, service.WrapAsInternalError(err)
	}
	if len(idGroup) == 0 {
		// 返回空集
		return []Account{}, nil
	}
	if err := tx.Model(&Account{}).Where("id in (?)", idGroup).Update(&Account{Status: StatusRunning}).Error; err != nil {
		return accounts, service.WrapAsInternalError(err)
	}
	if err := tx.Where("id in (?)", idGroup).Find(&accounts).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			// 返回空集
			return []Account{}, nil
		}
		return accounts, service.WrapAsInternalError(err)
	}
	return accounts, nil
}

func SetStatusNormal(db *database.DB, acc *Account) error {
	err := db.Model(acc).Update("status", StatusNormal).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	return nil
}
