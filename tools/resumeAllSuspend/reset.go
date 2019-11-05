package main

import (
	"github.com/inkedawn/JKWXRunner-server/service"
)

func main() {
	common := service.NewCommonService()
	common.Begin()
	defer func() {
		if x := recover(); x != nil {
			common.Rollback()
		}
	}()
	accSrv := service.NewAccountServiceUpon(common)
	err := accSrv.ResumeAllSuspend()
	if err != nil {
		panic(err)
	}
	common.Commit()
}
