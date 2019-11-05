package service

import (
	"database/sql"
	"errors"
	"sync"
	"time"

	ssmt "github.com/inkedawn/go-sunshinemotion/v3"
	"github.com/jinzhu/gorm"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv/accLogSrv"
)

var (
	ErrNoAccount           = errors.New("没有找到帐号")
	ErrAccountExistAlready = errors.New("帐号已存在")
)

type AccountStatus = string

//noinspection GoUnusedConst
const (
	AccountStatusNormal     AccountStatus = "normal"     // normal existence
	AccountStatusPause      AccountStatus = "pause"      // pause due to  human-reason, long-period
	AccountStatusRunning    AccountStatus = "running"    // running, can't be fetch by other executors
	AccountStatusFinished   AccountStatus = "finished"   // finished normally
	AccountStatusSuspend    AccountStatus = "suspend"    // suspend due to software error, short-period
	AccountStatusTerminated AccountStatus = "terminated" // processed completely, task is ready to be deleted
	AccountStatusAborted    AccountStatus = "aborted"    // aborted due to human-reason
	AccountStatusInQueue    AccountStatus = "inqueue"    // waitting to run, can't be fetch by other executors
)

type TaskRunResult = string

const (
	TaskRunSuccess       TaskRunResult = "success"
	TaskRunErrorOccurred TaskRunResult = "error"
)

//noinspection GoUnusedConst
const (
	UserRole_Normal = iota
	UserRole_Cheater
)

type IAccountService interface {
	CountAccounts() (n uint, err error)
	ListAccounts() ([]datamodels.Account, error)
	ListAccountsRange(offset, num uint) ([]datamodels.Account, error)
	SaveAccount(cc *datamodels.Account) error                                                    // Save update value in database, if the value doesn't have primary key(id), will insert it
	GetAccount(id uint) (*datamodels.Account, error)                                             // return ErrNoAccount if record not exist.
	GetAccountByStuNum(schoolID int64, stuNum string) (acc *datamodels.Account, err error)       // return ErrNoAccount if record not exist.
	GetActiveAccountByStuNum(schoolID int64, stuNum string) (acc *datamodels.Account, err error) // return ErrNoAccount if record not exist.
	SetCheckCheaterFlag(id uint, check bool) error
	CreateAccount(SchoolID int64, StuNum string, Password string) (*datamodels.Account, error)
	ResumeAllSuspend() error
}

type accountService struct {
	db *database.DB
	sync.Locker
}

func (a accountService) ResumeAllSuspend() error {
	a.Lock()
	defer a.Unlock()
	return a.db.Model(&datamodels.Account{}).
		Where("status=? AND last_result=?", AccountStatusSuspend, TaskRunErrorOccurred).
		Updates(&datamodels.Account{
			Status:   AccountStatusNormal,
			LastTime: sql.NullTime{Valid: false},
		}).Error
}

func (a accountService) GetActiveAccountByStuNum(schoolID int64, stuNum string) (acc *datamodels.Account, err error) {
	a.Lock()
	defer a.Unlock()
	acc = new(datamodels.Account)
	err = a.db.Where("school_id = ? AND stu_num = ? AND status IN (?)", schoolID, stuNum, []string{
		AccountStatusNormal,
		AccountStatusRunning,
		AccountStatusInQueue,
		AccountStatusFinished,
	}).Find(&acc).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNoAccount
	}
	if err != nil {
		return nil, WrapAsInternalError(err)
	}
	return acc, nil
}

func (a accountService) SetCheckCheaterFlag(id uint, check bool) error {
	a.Lock()
	defer a.Unlock()
	return a.setCheckCheaterFlagLocked(id, check)
}

func (a accountService) setCheckCheaterFlagLocked(id uint, check bool) error {
	return a.db.Model(&datamodels.Account{}).
		Select("check_cheat_marked").
		Where("id=?", id).
		Updates(map[string]interface{}{"check_cheat_marked": sql.NullBool{Valid: true, Bool: check}}).
		Error
}

func (a accountService) GetAccount(id uint) (acc *datamodels.Account, err error) {
	a.Lock()
	defer a.Unlock()
	return a.getAccountLocked(id)
}

func (a accountService) getAccountLocked(id uint) (acc *datamodels.Account, err error) {
	acc = new(datamodels.Account)
	err = a.db.Where("id=?", id).Find(&acc).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNoAccount
	}
	if err != nil {
		return nil, WrapAsInternalError(err)
	}
	return acc, nil
}

func (a accountService) GetAccountByStuNum(schoolID int64, stuNum string) (acc *datamodels.Account, err error) {
	a.Lock()
	defer a.Unlock()
	return a.getAccountByStuNumLocked(schoolID, stuNum)
}

func (a accountService) getAccountByStuNumLocked(schoolID int64, stuNum string) (acc *datamodels.Account, err error) {
	acc = new(datamodels.Account)
	err = a.db.Where("school_id = ? AND stu_num = ?", schoolID, stuNum).Find(&acc).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNoAccount
	}
	if err != nil {
		return nil, WrapAsInternalError(err)
	}
	return acc, nil
}

func (a accountService) CountAccounts() (n uint, err error) {
	a.Lock()
	defer a.Unlock()
	err = a.db.Model(&datamodels.Account{}).Count(&n).Error
	return
}

func (a accountService) ListAccounts() ([]datamodels.Account, error) {
	a.Lock()
	defer a.Unlock()
	var accounts []datamodels.Account
	if err := a.db.Find(&accounts).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			// 返回空集
			return accounts, nil
		}
		return accounts, WrapAsInternalError(err)
	}
	return accounts, nil
}

func (a accountService) ListAccountsRange(offset, num uint) ([]datamodels.Account, error) {
	a.Lock()
	defer a.Unlock()
	var accounts []datamodels.Account
	if err := a.db.Offset(offset).Limit(num).Find(&accounts).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			// 返回空集
			return accounts, nil
		}
		return accounts, WrapAsInternalError(err)
	}
	return accounts, nil
}

func (a accountService) SaveAccount(acc *datamodels.Account) error {
	a.Lock()
	defer a.Unlock()
	return a.saveAccountLocked(acc)
}

func (a accountService) saveAccountLocked(acc *datamodels.Account) error {
	newAcc := a.db.NewRecord(acc)
	err := a.db.Save(&acc).Error
	if err != nil {
		return WrapAsInternalError(err)
	}
	if newAcc {
		accLogSrv.AddLogSuccess(a.db, acc.ID, "创建成功")
	}
	return nil
}

func (a accountService) CreateAccount(SchoolID int64, StuNum string, Password string) (acc *datamodels.Account, err error) {
	a.Lock()
	defer a.Unlock()
	return a.createAccountLocked(SchoolID, StuNum, Password)
}

func (a accountService) createAccountLocked(SchoolID int64, StuNum string, Password string) (acc *datamodels.Account, err error) {
	tx := a.db.Begin()
	defer func() {
		// panic recovery
		if x := recover(); x != nil {
			tx.Rollback()
			err = x.(error)
		}
		// transaction finish
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	acc = &datamodels.Account{
		SchoolID: SchoolID,
		StuNum:   StuNum,
		Password: Password,
	}
	acc, err = a.getAccountByStuNumLocked(SchoolID, StuNum)
	switch err {
	case ErrNoAccount:
		break // okay, continue to create
	case nil:
		return acc, ErrAccountExistAlready
	default:
		panic(err)
	}
	ssmtDevice := ssmt.GenerateDevice()
	session := new(ssmt.Session)
	session.Device = ssmtDevice
	info, err := session.Login(SchoolID, StuNum, "123", ssmt.PasswordHash(Password))
	if err != nil {
		panic(err)
	}

	fetchTime := time.Now()
	sport, err := session.GetSportResult()
	if err != nil {
		panic(err)
	}
	err = NewUserSportResultServiceOn(tx).SaveCacheSportResult(datamodels.CacheUserSportResultFromSSMTSportResult(*sport, session.User.UserID, fetchTime))
	if err != nil {
		panic(err)
	}
	dev := datamodels.DeviceFromSSMTDevice(*ssmtDevice)
	err = NewDeviceServiceOn(tx).SaveDevice(&dev)
	if err != nil {
		panic(err)
	}

	const defaultOwnerID = 1
	limit := ssmt.GetDefaultLimitParams(info.Sex)
	acc = &datamodels.Account{
		OwnerID:          defaultOwnerID,
		SchoolID:         SchoolID,
		StuNum:           StuNum,
		Password:         Password,
		RunDistance:      limit.LimitTotalMaxDistance,
		DeviceID:         dev.ID,
		Status:           AccountStatusNormal,
		Memo:             "",
		CheckCheatMarked: sql.NullBool{Valid: false},
	}
	acc.RunDistance = ssmt.NormalizeDistance(acc.RunDistance)
	acc.StartDistance = sport.ActualDistance
	acc.FinishDistance = sport.QualifiedDistance

	err = NewAccountServiceOn(tx).SaveAccount(acc)
	if err != nil {
		panic(err)
	}
	err = NewUserIDRelServiceOn(tx).SaveRelation(acc.ID, session.User.UserID)
	if err != nil {
		panic(err)
	}
	if info.UserRoleID == UserRole_Cheater {
		err = a.setCheckCheaterFlagLocked(acc.ID, false)
		if err != nil {
		}
	}
	return acc, nil
}

func NewAccountService() IAccountService {
	return NewAccountServiceOn(database.GetDB())
}

func NewAccountServiceOn(db *database.DB) IAccountService {
	return &accountService{db: db, Locker: &sync.Mutex{}}
}
