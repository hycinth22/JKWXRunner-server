package main

import (
	"fmt"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/database/model"
)

func main() {
	db := database.GetDB()
	tx := db.Begin()
	tx.AutoMigrate(model.ModelsCollection...)
	errs := tx.GetErrors()
	if len(errs) != 0 {
		tx.Rollback()
		fmt.Println(errs)
		return
	}
	tx.Commit()
}
