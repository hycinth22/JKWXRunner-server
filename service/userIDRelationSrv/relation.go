package userIDRelationSrv

import (
	"errors"
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/database/model"
	"github.com/inkedawn/JKWXFucker-server/service"
)

var (
	ErrNotFound = errors.New("没有找到指定的UserID关系")
)

func GetLocalUID(db *database.DB, remoteUserID int64) (uint, error) {
	result := model.UserIDRelation{}
	err := db.First(&result, &model.UserIDRelation{RemoteUserID: remoteUserID}).Error
	if err != nil {
		if database.IsRecordNotFoundError(err) {
			return 0, ErrNotFound
		}
		return 0, service.WrapAsInternalError(err)
	}
	return result.UID, nil
}

// 保存UserInfo到缓存（通常在登录后保存）
func GetRemoteUserID(db *database.DB, uid uint) (int64, error) {
	result := model.UserIDRelation{}
	err := db.First(&result, &model.UserIDRelation{UID: uid}).Error
	if err != nil {
		if database.IsRecordNotFoundError(err) {
			return 0, ErrNotFound
		}
		return 0, service.WrapAsInternalError(err)
	}
	return result.RemoteUserID, nil
}

func SaveRelation(db *database.DB, localUID uint, remoteUserID int64) error {
	rel := model.UserIDRelation{
		UID:          localUID,
		RemoteUserID: remoteUserID,
	}
	err := db.Save(&rel).Error
	if err != nil {
		return service.WrapAsInternalError(err)
	}
	return nil
}
