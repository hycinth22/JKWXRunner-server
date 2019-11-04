package userIDRelationSrv

import (
	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service"
)

var (
	// DEPRECATED: use service.ErrUserIDRelNotFound
	ErrNotFound = service.ErrUserIDRelNotFound
)

// DEPRECATED: use service.IIUserIDRelService
func GetLocalUID(db *database.DB, remoteUserID int64) (uint, error) {
	return service.NewUserIDRelServiceOn(db).GetLocalUID(remoteUserID)
}

// DEPRECATED: use service.IIUserIDRelService
// 保存UserInfo到缓存（通常在登录后保存）
func GetRemoteUserID(db *database.DB, uid uint) (int64, error) {
	return service.NewUserIDRelServiceOn(db).GetRemoteUserID(uid)
}

// DEPRECATED: use service.IIUserIDRelService
func SaveRelation(db *database.DB, localUID uint, remoteUserID int64) error {
	return service.NewUserIDRelServiceOn(db).SaveRelation(localUID, remoteUserID)
}
