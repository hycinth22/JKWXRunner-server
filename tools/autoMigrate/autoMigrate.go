package main

import (
	"fmt"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
)

func main() {
	db := database.GetDB()
	tx := db.Begin()
	tx.AutoMigrate(datamodels.ModelsCollection...)
	errs := tx.GetErrors()
	if len(errs) != 0 {
		tx.Rollback()
		fmt.Println(errs)
		return
	}
	tx.Commit()
}
