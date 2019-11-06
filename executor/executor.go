package main

import (
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv"
)

const libVersion = ssmt.AppVersionID

func main() {
	// TODO: Watch Signal
	if !VersionCheck() {
		os.Exit(1)
	}
	db := database.GetDB()
	accounts, err := accountSrv.ListAndSetRunStatusForAllAccountsWaitRun(db)
	if err != nil {
		panic(err)
	}
	log.Println("Run List:")
	for _, acc := range accounts {
		log.Println(dumpStructValue(acc))
	}
	rand.Shuffle(len(accounts), func(i, j int) {
		accounts[i], accounts[j] = accounts[j], accounts[i]
	})
	const (
		nWorker    = 3
		retryTimes = 3
	)
	var accAllGroups [nWorker][]*datamodels.Account
	for i, acc := range accounts {
		target := i % nWorker
		accAllGroups[target] = append(accAllGroups[target], acc)
	}
	var wg sync.WaitGroup
	for _, accGroup := range accAllGroups {
		startupTaskWorker(service.NewCommonService(), accGroup, &wg, retryTimes)
		wg.Add(1)
	}
	wg.Wait()
	log.Println("All Account Task has been completed")
}

func VersionCheck() bool {
	log.Println("Target App version: ", ssmt.AppVersionID)
	s := ssmt.CreateSession()
	info, err := s.GetAppInfo()
	if err != nil {
		log.Println("Can't get latest app info.", err)
		return false
	}
	if info.VerNumber > libVersion {
		log.Println("Latest App version: ", info.VerNumber)
		log.Println("Need to upgrade!!!")
		return false
	}
	return true
}
