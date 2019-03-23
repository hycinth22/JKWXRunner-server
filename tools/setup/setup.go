package main

import (
	"fmt"
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/database/model"
)

func main() {
	db := database.GetDB()
	tx := db.Begin()
	tx.CreateTable(&model.Account{}, &model.AccountLog{})
	tx.CreateTable(&model.Token{}, &model.Device{})
	tx.CreateTable(&model.UserIDRelation{}, &model.CacheUserInfo{}, &model.CacheUserSportResult{})

	errs := tx.GetErrors()
	if len(errs) != 0 {
		tx.Rollback()
		fmt.Println(errs)
		return
	}
	tx.Commit()
}
