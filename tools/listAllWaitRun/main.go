package main

import (
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/service/accountSrv"
	"fmt"
)

func main() {
	db := database.GetDB()

	accounts, err := accountSrv.ListAllAccountsWaitRun(db)
	if err != nil {
		panic(err)
	}
	for i := range accounts {
		fmt.Println(accounts[i])
	}
}