package model

import (
	"time"
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
