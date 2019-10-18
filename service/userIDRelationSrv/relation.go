package userIDRelationSrv

import (
	"errors"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service"
)

var (
	ErrNotFound = errors.New("没有找到指定的UserID关系")
)

func GetLocalUID(db *database.DB, remoteUserID int64) (uint, error) {
	return service.NewUserIDRelServiceOn(db).GetLocalUID(remoteUserID)
}

// 保存UserInfo到缓存（通常在登录后保存）
func GetRemoteUserID(db *database.DB, uid uint) (int64, error) {
	return service.NewUserIDRelServiceOn(db).GetRemoteUserID(uid)
}

func SaveRelation(db *database.DB, localUID uint, remoteUserID int64) error {
	return service.NewUserIDRelServiceOn(db).SaveRelation(localUID, remoteUserID)
}
