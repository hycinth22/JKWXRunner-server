package datamodels

import (
	"time"

	ssmt "github.com/inkedawn/go-sunshinemotion/v3"
)

// 缓存的UserInfo信息（通常是上次登录时保存的）
// 对于每个UserID只保存对应唯一的一份
type CacheUserInfo struct {
	RemoteUserID int64     `gorm:"primary_key;NOT NULL"`
	FetchTime    time.Time `gorm:"NOT NULL"`

	ClassID       int64  `gorm:"NOT NULL"`
	ClassName     string `gorm:"NOT NULL"`
	CollegeID     int64  `gorm:"NOT NULL"`
	CollegeName   string `gorm:"NOT NULL"`
	SchoolID      int64  `gorm:"NOT NULL"`
	SchoolName    string `gorm:"NOT NULL"`
	SchoolNumber  string `gorm:"NOT NULL"`
	NickName      string `gorm:"NOT NULL"`
	StudentName   string `gorm:"NOT NULL"`
	StudentNumber string `gorm:"NOT NULL"`
	IsTeacher     int    `gorm:"NOT NULL"`
	Sex           string `gorm:"NOT NULL"`
	PhoneNumber   string `gorm:"NOT NULL"`
	UserRoleID    int    `gorm:"NOT NULL"`
}

func CacheUserInfoFromSSMTUserInfo(info ssmt.UserInfo, userID int64, fetchTime time.Time) CacheUserInfo {
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
