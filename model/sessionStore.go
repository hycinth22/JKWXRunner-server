package model

import (
	sunshinemotion "../sunshinemotion"
	"bytes"
	"encoding/gob"
	"errors"
)

type SessionStore struct {
	AccountID  uint       `gorm:"primary_key"`
	SessionObj []byte
}

var ErrSessionNotFound = errors.New("sessionStore not found")

func SaveSession(accountID uint, session *sunshinemotion.Session) (err error) {
	buffer := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buffer)
	if err := enc.Encode(session); err != nil {
		return errors.New("saveSession encode Fail" + err.Error())
	}

	if err := db.Save(&SessionStore{
		AccountID:  accountID,
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
		return nil, errors.New("getSession query sesion Fail" + err.Error())
	}
	if db.RecordNotFound() {
		return nil, ErrSessionNotFound
	}
	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(store.SessionObj)
	dec := gob.NewDecoder(buffer)
	session = new(sunshinemotion.Session)
	if err := dec.Decode(&session); err != nil {
		return nil, errors.New("getSession decode Fail" + err.Error())
	}
	return session, nil
}