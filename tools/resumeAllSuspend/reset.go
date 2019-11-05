package main

import (
	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service"
)

func main() {
	tx := database.GetDB().Begin()
	defer func() {
		if x := recover(); x != nil {
			tx.Rollback()
		}
	}()
	accSrv := service.NewAccountServiceOn(tx)
	err := accSrv.ResumeAllSuspend()
	if err != nil {
		panic(err)
	}
	tx.Commit()
}
