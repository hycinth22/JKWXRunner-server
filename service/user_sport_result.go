package service

import (
	"errors"
	"sync"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
)

var (
	ErrNoUserSportResult = errors.New("没有找到该用户缓存的运动结果")
)

type IUserSportResultService interface {
	GetCacheSportResult(userID int64) (datamodels.CacheUserSportResult, error)
	SaveCacheSportResult(info datamodels.CacheUserSportResult) error
	GetLocalUserCacheSportResult(localUID uint) (info datamodels.CacheUserSportResult, err error)
}

type userSportResultSrv struct {
	db *database.DB
	sync.Locker
}

// 从数据库获取缓存的信息（通常是上次执行运动任务时更新的）
func (u *userSportResultSrv) GetCacheSportResult(userID int64) (datamodels.CacheUserSportResult, error) {
	// TODO: lock
	var info datamodels.CacheUserSportResult
	if err := u.db.First(&info, &datamodels.CacheUserSportResult{RemoteUserID: userID}).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			return info, ErrNoUserSportResult
		}
		return info, WrapAsInternalError(err)
	}
	return info, nil
}

// 保存SportResult到缓存（通常是上次执行运动任务时更新的）
func (u *userSportResultSrv) SaveCacheSportResult(info datamodels.CacheUserSportResult) error {
	// TODO: lock
	err := u.db.Save(&info).Error
	if err != nil {
		return WrapAsInternalError(err)
	}
	return nil
}

func (u *userSportResultSrv) GetLocalUserCacheSportResult(localUID uint) (info datamodels.CacheUserSportResult, err error) {
	// TODO: lock
	remoteUID, err := NewUserIDRelServiceOn(u.db).GetRemoteUserID(localUID)
	if err == ErrUserIDRelNotFound {
		err = ErrNoUserSportResult
		return
	} else if err != nil {
		return
	}
	info, err = u.GetCacheSportResult(remoteUID)
	return
}

func NewUserSportResultService() IUserSportResultService {
	return NewUserSportResultServiceOn(database.GetDB())
}
func NewUserSportResultServiceOn(db *database.DB) IUserSportResultService {
	return &userSportResultSrv{db: db}
}
