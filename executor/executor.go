package main

import (
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv"
)

const libVersion = ssmt.AppVersionID

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	const retryTimes = 3
	var wg sync.WaitGroup
	nWorker := len(accounts)
	const (
		eachAroundTime = 3 * time.Minute
		timeLimit      = 6 * time.Hour
	)
	totalTime := time.Duration(nWorker) * eachAroundTime
	if totalTime > timeLimit {
		totalTime = timeLimit
	}
	for i, acc := range accounts {
		if i != 0 {
			totalTime -= sleepPartOfTotalTime(int64(nWorker), totalTime)
		}
		accGroup := []*datamodels.Account{acc}
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
