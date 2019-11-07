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
// 保存UserInfo到缓存（通常在登录后保存）
func GetRemoteUserID(db *database.DB, uid uint) (int64, error) {
	return service.NewUserIDRelServiceOn(db).GetRemoteUserID(uid)
}
