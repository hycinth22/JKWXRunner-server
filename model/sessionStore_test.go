package model

import (
	sunshinemotion "../sunshinemotion"
	"log"
	"testing"
)

func TestSaveSession(t *testing.T) {
	session := sunshinemotion.CreateSession()
	if err := session.Login("021840104", "123", sunshinemotion.PasswordHash("123456")); err != nil {
		t.FailNow()
	}
	if err := SaveSession(111111111, session); err != nil {
		t.FailNow()
	}
}
func TestGetSession(t *testing.T) {
	session, err := GetSession(111111111)
	if err != nil {
		t.Fail()
		return
	}
	log.Println(session)
}
