package service

import (
	"errors"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
)

var (
	ErrUserIDRelNotFound = errors.New("没有找到帐号")
)

type IUserIDRelService interface {
	GetLocalUID(remoteUserID int64) (uint, error)
	GetRemoteUserID(uid uint) (int64, error)
	SaveRelation(localUID uint, remoteUserID int64) error
}

type userIDRelService struct {
	db *database.DB
}

func (u userIDRelService) GetLocalUID(remoteUserID int64) (uint, error) {
	result := datamodels.UserIDRelation{}
	err := u.db.First(&result, &datamodels.UserIDRelation{RemoteUserID: remoteUserID}).Error
	if err != nil {
		if database.IsRecordNotFoundError(err) {
			return 0, ErrUserIDRelNotFound
		}
		return 0, WrapAsInternalError(err)
	}
	return result.UID, nil
}

func (u userIDRelService) GetRemoteUserID(uid uint) (int64, error) {
	result := datamodels.UserIDRelation{}
	err := u.db.First(&result, &datamodels.UserIDRelation{UID: uid}).Error
	if err != nil {
		if database.IsRecordNotFoundError(err) {
			return 0, ErrUserIDRelNotFound
		}
		return 0, WrapAsInternalError(err)
	}
	return result.RemoteUserID, nil
}

func (u userIDRelService) SaveRelation(localUID uint, remoteUserID int64) error {
	rel := datamodels.UserIDRelation{
		UID:          localUID,
		RemoteUserID: remoteUserID,
	}
	err := u.db.Save(&rel).Error
	if err != nil {
		return WrapAsInternalError(err)
	}
	return nil
}

func NewUserIDRelService() IUserIDRelService {
	return &userIDRelService{db: database.GetDB()}
}
func NewUserIDRelServiceOn(db *database.DB) IUserIDRelService {
	return &userIDRelService{db: db}
}
