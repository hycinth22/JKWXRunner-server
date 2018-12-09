package main

import (
	"github.com/inkedawn/JKWXFucker-server/model"
	sunshinemotion "github.com/inkedawn/go-sunshinemotion"
	"log"
)

func main() {
	list := model.ListStoredSessionAccounts()
	log.Println("find ", len(list), "account")
	for _, accountID := range list {
		s, err := model.GetSession(accountID)
		if err != nil {
			log.Println(accountID, err)
			continue
		}
		if s.PhoneIMEI == "" {
			s.PhoneIMEI = sunshinemotion.GenerateIMEI()
		}
		if s.PhoneModel == "" {
			s.PhoneModel = sunshinemotion.RandModel()
		}
		err = nil
		model.SaveSession(accountID, s)
		if err == nil {
			log.Println(accountID, "save successfully.")
		} else {
			log.Println(accountID, "save failed, err ", err.Error())
		}
	}
}
