package main

import (
	"fmt"
	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv"
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
