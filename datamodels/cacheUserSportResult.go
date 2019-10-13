package datamodels

import (
	"time"
)

// 缓存的UserSportResult信息（通常是发起获取SportResult请求时保存的）
// 对于每个UserID只保存对应唯一的一份
type CacheUserSportResult struct {
	RemoteUserID int64     `gorm:"primary_key;NOT NULL"`
	FetchTime    time.Time `gorm:"NOT NULL"`

	Year              int       `gorm:"NOT NULL"` // 年度
	Term              string    `gorm:"NOT NULL"` // 学期
	QualifiedDistance float64   `gorm:"NOT NULL"` // 达标距离
	ComputedDistance  float64   `gorm:"NOT NULL"` // 已计距离
	LastTime          time.Time `gorm:"NOT NULL"` // 上次跑步时间
}
