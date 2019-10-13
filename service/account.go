package service

import (
	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
)

type IAccountService interface {
	CountAccounts() (n uint, err error)
	ListAccounts() ([]datamodels.Account, error)
	ListAccountsRange(offset, num uint) ([]datamodels.Account, error)
	SaveAccount(cc *datamodels.Account) // Save update value in database, if the value doesn't have primary key(id), will insert it
}

type accountService struct {
	db *database.DB
}

func (a accountService) CountAccounts() (n uint, err error) {
	panic("implement me")
}

func (a accountService) ListAccounts() ([]datamodels.Account, error) {
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
	panic("implement me")
}

func (a accountService) SaveAccount(cc *datamodels.Account) {
	panic("implement me")
}

func NewAccountService() IAccountService {
	return &accountService{db: database.GetDB()}
}
