// 提供对Account的管理
package accountSrv

import (
	"errors"
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/database/model"
	"github.com/inkedawn/JKWXFucker-server/service"
	"github.com/inkedawn/JKWXFucker-server/service/accountSrv/accLogSrv"

	"time"
)

type Account = model.Account

var (
	ErrNoAccount = errors.New("没有找到帐号")
)

type Status = string

//noinspection GoUnusedConst
const (
	StatusNormal     Status = "normal"     // normal existence
	StatusPause      Status = "pause"      // pause due to  human-reason, long-period
	StatusRunning    Status = "running"    // running, can't be fetch by other executors
	StatusFinished   Status = "finished"   // finished normally
	StatusSuspend    Status = "suspend"    // suspend due to software error, short-period
	StatusTerminated Status = "terminated" // processed completely, task is ready to be deleted
	StatusAborted    Status = "aborted"    // aborted due to human-reason
)

type RunResult = string

const (
	RunSuccess       RunResult = "success"
	RunErrorOccurred RunResult = "error"
)

func CountAccounts(db *database.DB) (n uint, err error) {
	err = db.Model(&Account{}).Count(&n).Error
	return
}

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
	newAcc := db.NewRecord(acc)
	err := db.Save(&acc).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	if newAcc {
		accLogSrv.AddLogSuccess(db, acc.ID, "创建成功")
	}
	return nil
}

func ListAllAccountsWaitRun(db *database.DB) (accounts []Account, err error) {
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
	if err := tx.Where("id in (?)", idGroup).Find(&accounts).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			// 返回空集
			return []Account{}, nil
		}
		return accounts, service.WrapAsInternalError(err)
	}
	return accounts, nil
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

func SetStatus(db *database.DB, acc *Account, status Status) error {
	err := db.Model(acc).Update("status", status).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	acc.Status = status
	return nil
}

func SetStatusNormal(db *database.DB, acc *Account) error {
	return SetStatus(db, acc, StatusNormal)
}

func SetLastTime(db *database.DB, acc *Account, lastTime time.Time) error {
	err := db.Model(acc).Update("last_time", lastTime).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	acc.LastTime = lastTime
	return nil
}

func SetLastResult(db *database.DB, acc *Account, r RunResult) error {
	err := db.Model(acc).Update("last_result", r).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	acc.LastResult = r
	return nil
}
