package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/inkedawn/JKWXFucker-server/model"
	sunshinemotion "github.com/inkedawn/go-sunshinemotion"
	"log"
)

func main() {
	var list []model.SessionStore
	db := model.GetDB()
	db.AutoMigrate(&model.SessionStore{})
	if err := db.Raw("SELECT DISTINCT * FROM session_store").Scan(&list).Error; err != nil {
		panic(err)
	}
	for _, a := range list {
		s, err := buildSessionObj(a.SessionObj)
		if err != nil {
			log.Println(a, err)
			continue
		}
		a.Username = s.UserInfo.StudentNumber
		err = db.Save(a).Error
		if err != nil {
			log.Println(a, err)
		}
	}
}

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
