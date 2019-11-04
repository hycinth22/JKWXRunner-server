package datamodels

type UserIDRelation struct {
	UID          uint  `gorm:"primary_key"`
	RemoteUserID int64 `gorm:"index;NOT NULL"`
}
