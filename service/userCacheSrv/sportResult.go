package userCacheSrv

import (
	"errors"
	"time"

	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/database/model"
	"github.com/inkedawn/JKWXRunner-server/service"
)

var (
	ErrNoSportResult = errors.New("没有找到该用户缓存的运动结果")
)

type CacheSportResult = model.CacheUserSportResult

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

// 保存SportResult到缓存（通常是上次执行运动任务时更新的）
func SaveCacheSportResult(db *database.DB, info CacheSportResult) error {
	err := db.Save(&info).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	return nil
}

func FromSSMTSportResult(info ssmt.SportResult, userID int64, fetchTime time.Time) CacheSportResult {
	return CacheSportResult{
		RemoteUserID:      userID,
		FetchTime:         fetchTime,
		Year:              info.Year,
		Term:              info.Term,
		QualifiedDistance: info.QualifiedDistance,
		ComputedDistance:  info.ActualDistance,
		LastTime:          info.LastTime,
	}
}
