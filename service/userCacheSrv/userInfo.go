// 提供对远程用户数据缓存的存取支持
package userCacheSrv

import (
	"errors"
	"time"

	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
)

var (
	// DEPRECATED: use service.ErrNoUserInfo
	ErrNoUserInfo = errors.New("没有找到该用户缓存的用户信息")
)

// DEPRECATED: use datamodels.CacheUserInfo
type CacheUserInfo = datamodels.CacheUserInfo

//noinspection GoUnusedConst
const (
	// DEPRECATED: use service.UserRole_Normal
	UserRole_Normal = service.UserRole_Normal
	// DEPRECATED: use service.UserRole_Cheater
	UserRole_Cheater = service.UserRole_Cheater
)

// 从数据库获取缓存的信息（通常是上次登录时保存的）
func GetCacheUserInfo(db *database.DB, userID int64) (CacheUserInfo, error) {
	var info CacheUserInfo
	if err := db.First(&info, &CacheUserInfo{RemoteUserID: userID}).Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			return info, ErrNoUserInfo
		}
		return info, service.WrapAsInternalError(err)
	}
	return info, nil
}

// 保存UserInfo到缓存（通常在登录后保存）
func SaveCacheUserInfo(db database.TX, info CacheUserInfo) error {
	err := db.Save(&info).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	return nil
}

// DEPRECATED: use datamodels.CacheUserInfoFromSSMTUserInfo
func FromSSMTUserInfo(info ssmt.UserInfo, userID int64, fetchTime time.Time) CacheUserInfo {
	return datamodels.CacheUserInfoFromSSMTUserInfo(info, userID, fetchTime)
}
