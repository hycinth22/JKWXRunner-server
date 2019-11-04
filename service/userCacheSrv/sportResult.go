package userCacheSrv

import (
	"time"

	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
	"github.com/inkedawn/JKWXRunner-server/service/userIDRelationSrv"
)

var (
	// DEPRECATED: use service.ErrNoUserSportResult
	ErrNoSportResult = service.ErrNoUserSportResult
)

// DEPRECATED: use datamodels.CacheUserSportResult
type CacheSportResult = datamodels.CacheUserSportResult

// DEPRECATED: use service.IUserSportResultService
// 从数据库获取缓存的信息（通常是上次执行运动任务时更新的）
func GetCacheSportResult(db *database.DB, userID int64) (CacheSportResult, error) {
	var info CacheSportResult
	if err := db.First(&info, &CacheSportResult{RemoteUserID: userID}).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			return info, ErrNoSportResult
		}
		return info, service.WrapAsInternalError(err)
	}
	return info, nil
}

// DEPRECATED: use service.IUserSportResultService
// 保存SportResult到缓存（通常是上次执行运动任务时更新的）
func SaveCacheSportResult(db *database.DB, info CacheSportResult) error {
	err := db.Save(&info).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	return nil
}

// DEPRECATED: use datamodels.CacheUserSportResultFromSSMTSportResult
func FromSSMTSportResult(info ssmt.SportResult, userID int64, fetchTime time.Time) CacheSportResult {
	return datamodels.CacheUserSportResultFromSSMTSportResult(info, userID, fetchTime)
}

// DEPRECATED: use service.IUserSportResultService
func GetLocalUserCacheSportResult(db *database.DB, localUID uint) (info CacheSportResult, err error) {
	remoteUID, err := userIDRelationSrv.GetRemoteUserID(db, localUID)
	if err == userIDRelationSrv.ErrNotFound {
		err = ErrNoSportResult
		return
	} else if err != nil {
		return
	}
	info, err = GetCacheSportResult(db, remoteUID)
	return
}
