package service

import (
	"errors"
	"sync"

	"github.com/jinzhu/gorm"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv/accLogSrv"
)

var (
	ErrNoAccount = errors.New("没有找到帐号")
)

type IAccountService interface {
	CountAccounts() (n uint, err error)
	ListAccounts() ([]datamodels.Account, error)
	ListAccountsRange(offset, num uint) ([]datamodels.Account, error)
	SaveAccount(cc *datamodels.Account) error                                              // Save update value in database, if the value doesn't have primary key(id), will insert it
	GetAccount(id uint) (*datamodels.Account, error)                                       // return ErrNoAccount if record not exist.
	GetAccountByStuNum(schoolID int64, stuNum string) (acc *datamodels.Account, err error) // return ErrNoAccount if record not exist.
	SetCheckCheaterFlag(id uint, check bool) error
}

type accountService struct {
	db *database.DB
	sync.Locker
}

func (a accountService) SetCheckCheaterFlag(id uint, check bool) error {
	a.Lock()
	defer a.Unlock()
	a.db.Model(&datamodels.Account{}).Select("check_cheat_marked").Updates(map[string]interface{}{"check_cheat_marked": false})
	return a.db.Error
}

func (a accountService) GetAccount(id uint) (acc *datamodels.Account, err error) {
	a.Lock()
	defer a.Unlock()
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

func NewAccountService() IAccountService {
	return NewAccountServiceOn(database.GetDB())
}

func NewAccountServiceOn(db *database.DB) IAccountService {
	return &accountService{db: db, Locker: &sync.Mutex{}}
}
