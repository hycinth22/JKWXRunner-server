package sessionSrv

import (
	"errors"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"

	"github.com/inkedawn/go-sunshinemotion/v3"
)

type token = datamodels.Token

var (
	ErrNoToken = errors.New("没有找到该用户的Token")
)

func getToken(db *database.DB, remoteUserID int64) (token, error) {
	var userToken token
	err := db.First(&userToken, &token{RemoteUserID: remoteUserID}).Error
	if err != nil && database.IsRecordNotFoundError(err) {
		return userToken, ErrNoToken
	}
	return userToken, err
}
func getTokenByUID(db *database.DB, uid uint) (t token, err error) {
	var userID datamodels.UserIDRelation
	err = db.First(&userID, datamodels.UserIDRelation{
		UID: uid,
	}).Error
	if err != nil {
		if database.IsRecordNotFoundError(err) {
			err = ErrNoToken
			return
		}
		return
	}
	return getToken(db, userID.RemoteUserID)
}

func saveToken(db *database.DB, userToken token) (err error) {
	err = db.Save(&userToken).Error
	if err != nil {
		return err
	}
	return nil
}

func fromSSMTToken(userID int64, userToken ssmt.UserToken) token {
	return token{
		RemoteUserID:   userID,
		TokenID:        userToken.TokenID,
		ExpirationTime: userToken.ExpirationTime,
	}
}

func toSSMTToken(t token) ssmt.UserToken {
	return ssmt.UserToken{
		TokenID:        t.TokenID,
		ExpirationTime: t.ExpirationTime,
	}
}
