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
	ErrNoUserInfo = errors.New("没有找到该用户缓存的用户信息")
)

type CacheUserInfo = datamodels.CacheUserInfo

//noinspection GoUnusedConst
const (
	UserRole_Normal = iota
	UserRole_Cheater
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
func SaveCacheUserInfo(db *database.DB, info CacheUserInfo) error {
	err := db.Save(&info).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	return nil
}

func FromSSMTUserInfo(info ssmt.UserInfo, userID int64, fetchTime time.Time) CacheUserInfo {
	return CacheUserInfo{
		RemoteUserID:  userID,
		FetchTime:     fetchTime,
		ClassID:       info.ClassID,
		ClassName:     info.ClassName,
		CollegeID:     info.CollegeID,
		CollegeName:   info.CollegeName,
		SchoolID:      info.SchoolID,
		SchoolName:    info.SchoolName,
		SchoolNumber:  info.SchoolNumber,
		NickName:      info.NickName,
		StudentName:   info.StudentName,
		StudentNumber: info.StudentNumber,
		IsTeacher:     info.IsTeacher,
		Sex:           info.Sex,
		PhoneNumber:   info.PhoneNumber,
		UserRoleID:    info.UserRoleID,
	}
}
