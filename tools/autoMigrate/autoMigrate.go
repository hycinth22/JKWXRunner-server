package main

import (
	"fmt"
	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/database/model"
)

func main() {
	db := database.GetDB()
	tx := db.Begin()
	tx.AutoMigrate(&model.Account{}, &model.AccountLog{})
	tx.AutoMigrate(&model.Token{}, &model.Device{})
	tx.AutoMigrate(&model.UserIDRelation{}, &model.CacheUserInfo{}, &model.CacheUserSportResult{})

	errs := tx.GetErrors()
	if len(errs) != 0 {
		tx.Rollback()
		fmt.Println(errs)
		return
	}
	tx.Commit()
}
