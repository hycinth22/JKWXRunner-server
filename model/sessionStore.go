package model

import (
	"bytes"
	"encoding/gob"
	"errors"
	sunshinemotion "github.com/inkedawn/go-sunshinemotion"
	"github.com/jinzhu/gorm"
	"log"
)

type SessionStore struct {
	AccountID  uint   `gorm:"primary_key"`
	Username   string `gorm:"index"`
	SessionObj []byte
}

var ErrSessionNotFound = errors.New("sessionStore not found")

func buildSessionObj(bin []byte) (*sunshinemotion.Session, error) {
	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(bin)
	dec := gob.NewDecoder(buffer)
	session := new(sunshinemotion.Session)
	if err := dec.Decode(&session); err != nil {
		return nil, errors.New("buildSessionObj decode Fail" + err.Error())
	}
	return session, nil
}

func SaveSession(accountID uint, session *sunshinemotion.Session) (err error) {
	log.Println("SaveSession, user", session.UserInfo.StudentNumber, session)

	buffer := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buffer)
	if err := enc.Encode(session); err != nil {
		return errors.New("saveSession encode Fail" + err.Error())
	}

	if err := db.Save(&SessionStore{
		AccountID:  accountID,
		Username:   session.UserInfo.StudentNumber,
		SessionObj: buffer.Bytes(),
	}).Error; err != nil {
		return errors.New("saveSession Fail" + err.Error())
	}
	return nil
}

func GetSession(accountID uint) (session *sunshinemotion.Session, err error) {
	store := &SessionStore{
		AccountID: accountID,
	}
	if err := db.First(&store).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrSessionNotFound
		}
		return nil, errors.New("getSession query sesion Fail" + err.Error())
	}
	session, err = buildSessionObj(store.SessionObj)
	log.Println("GetSession, user", session.UserInfo.StudentNumber, session)
	return session, nil
}

func GetSessionByUsername(username string) (session *sunshinemotion.Session, err error) {
	store := &SessionStore{}
	if err := db.Where("username = ?", username).First(&store).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrSessionNotFound
		}
		return nil, errors.New("getSession query sesion Fail" + err.Error())
	}
	session, err = buildSessionObj(store.SessionObj)
	log.Println("GetSession, user", session.UserInfo.StudentNumber, session)
	return session, nil
}

func ListStoredSessionAccounts() (allAccountID []uint) {
	var list []struct {
		AccountID uint
	}
	if err := db.Raw("SELECT DISTINCT account_id FROM session_store").Scan(&list).Error; err != nil {
		return nil
	}
	for _, a := range list {
		allAccountID = append(allAccountID, a.AccountID)
	}
	return allAccountID
}
