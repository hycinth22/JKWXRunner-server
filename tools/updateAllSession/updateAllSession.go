package main

import (
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/service/accountSrv"
	"github.com/inkedawn/JKWXFucker-server/service/sessionSrv"
)

func main() {
	db := database.GetDB()

	accounts, err := accountSrv.ListAccounts(db, 0, 10000000)
	if err != nil {
		panic(err)
	}
	for i := range accounts {
		println(i)
		if accounts[i].Status == accountSrv.StatusNormal {
			err := sessionSrv.UpdateSession(db, accounts[i])
			if err != nil {
				println(err.Error())
			}
		}
	}
}
