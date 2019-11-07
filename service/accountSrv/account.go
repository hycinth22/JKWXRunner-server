// 提供对Account的管理
package accountSrv

import (
	"database/sql"
	"time"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
)

// DEPRECATED
type Account = datamodels.Account // DEPRECATED

// DEPRECATED
var (
	ErrNoAccount = service.ErrNoAccount
)

// DEPRECATED
type Status = service.AccountStatus

// DEPRECATED
//noinspection GoUnusedConst
const (
	StatusNormal     Status = "normal"     // normal existence.
	StatusPause      Status = "pause"      // pause due to  human-reason, long-period
	StatusRunning    Status = "running"    // running, can't be fetch by other executors
	StatusFinished   Status = "finished"   // finished normally
	StatusSuspend    Status = "suspend"    // suspend due to software error, short-period
	StatusTerminated Status = "terminated" // processed completely, task is ready to be deleted
	StatusAborted    Status = "aborted"    // aborted due to human-reason
	StatusInQueue    Status = "inqueue"    // waitting to run, can't be fetch by other executors
)

// DEPRECATED
type RunResult = service.TaskRunResult

// DEPRECATED
func ListAccounts(db *database.DB, offset, num uint) ([]Account, error) {
	return service.NewAccountServiceOn(db).ListAccountsRange(offset, num)
}

// DEPRECATED
// Save update value in database, if the value doesn't have primary key(id), will insert it
func SaveAccount(db *database.DB, acc *Account) error {
	return service.NewAccountServiceOn(db).SaveAccount(acc)
}

// DEPRECATED
// return ErrNoAccount if record not exist.
func GetAccount(db *database.DB, id uint) (*Account, error) {
	return service.NewAccountServiceOn(db).GetAccount(id)
}

func ListAllAccountsWaitRun(db *database.DB) (accounts []datamodels.Account, err error) {
	now := time.Now()
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	common := service.NewCommonServiceOn(db)
	tx := common.Begin()
	defer func() {
		if err == nil {
			common.Commit()
		} else {
			common.Rollback()
		}
	}()
	var idGroup []uint
	if err := tx.Model(&datamodels.Account{}).Where("status = ? AND (last_time < ? OR last_time IS NULL)", StatusNormal, todayZero).Pluck("id", &idGroup).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			// 返回空集
			return []datamodels.Account{}, nil
		}
		return accounts, service.WrapAsInternalError(err)
	}
	if len(idGroup) == 0 {
		// 返回空集
		return []datamodels.Account{}, nil
	}
	if err := tx.Where("id in (?)", idGroup).Find(&accounts).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			// 返回空集
			return []datamodels.Account{}, nil
		}
		return accounts, service.WrapAsInternalError(err)
	}
	return accounts, nil
}
func ListAndSetRunStatusForAllAccountsWaitRun(common service.ICommonService) (accounts []*datamodels.Account, err error) {
	now := time.Now()
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	tx := common.Begin()
	defer func() {
		if err == nil {
			common.Commit()
		} else {
			common.Rollback()
		}
	}()
	var idGroup []uint
	if err := tx.Model(&datamodels.Account{}).Where("status = ? AND (last_time < ? OR last_time IS NULL)", service.AccountStatusNormal, todayZero).Pluck("id", &idGroup).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			// 返回空集
			return []*datamodels.Account{}, nil
		}
		return accounts, service.WrapAsInternalError(err)
	}
	if len(idGroup) == 0 {
		// 返回空集
		return []*datamodels.Account{}, nil
	}
	if err := tx.Model(&datamodels.Account{}).Where("id in (?)", idGroup).Update(&datamodels.Account{Status: service.AccountStatusInQueue}).Error; err != nil {
		return accounts, service.WrapAsInternalError(err)
	}
	if err := tx.Where("id in (?)", idGroup).Find(&accounts).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			// 返回空集
			return []*datamodels.Account{}, nil
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

func SetLastTime(db *database.DB, acc *Account, lastTime time.Time) error {
	err := db.Model(acc).Update("last_time", lastTime).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	acc.LastTime = sql.NullTime{Valid: true, Time: lastTime}
	return nil
}

func SetLastResult(db *database.DB, acc *Account, r RunResult) error {
	err := db.Model(acc).Update("last_result", r).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	acc.LastResult = sql.NullString{Valid: true, String: r}
	return nil
}
