package main

import (
	"github.com/inkedawn/JKWXFucker-server/model"
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
		s.UpdateLimitParams()
		model.SaveSession(accountID, s)
		if err == nil {
			log.Println(accountID, "save successfully.")
		} else {
			log.Println(accountID, "save failed, err ", err.Error())
		}
	}
}
