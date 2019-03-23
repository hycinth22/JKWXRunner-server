package model

type UserIDRelation struct {
	UID          uint  `gorm:"primary_key"`
	RemoteUserID int64 `gorm:"UNIQUE_INDEX;NOT NULL"`
}
