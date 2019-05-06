package main

import (
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/service/accountSrv"
	"github.com/inkedawn/JKWXFucker-server/service/sessionSrv"
	"github.com/inkedawn/JKWXFucker-server/utils"
	"time"
)

func main() {
	db := database.GetDB()

	accounts, err := accountSrv.ListAccounts(db, 0, 10000000)
	if err != nil {
		panic(err)
	}
	n := len(accounts)
	for i := range accounts[1:] {
		println(i)
		if accounts[i].Status == accountSrv.StatusNormal {
			err := sessionSrv.UpdateSession(db, accounts[i])
			if err != nil {
				println(err.Error())
			}
			utils.SleepPartOfTotalTime(n, time.Duration(n)*5*time.Second)
		}
	}
}
