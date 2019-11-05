package main

import (
	"fmt"

	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
)

func main() {
	common := service.NewCommonService()
	tx := common.Begin()
	tx.AutoMigrate(datamodels.ModelsCollection...)
	errs := tx.GetErrors()
	if len(errs) != 0 {
		common.Rollback()
		fmt.Println(errs)
		return
	}
	common.Commit()
}
