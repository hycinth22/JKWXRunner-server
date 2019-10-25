package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv"
)

const lib_version = ssmt.AppVersionID

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
	const retryTimes = 3
	var wg sync.WaitGroup
	wg.Add(len(accounts))
	if len(accounts) >= 1 {
		startupTaskWorker(db, &accounts[0], &wg, retryTimes)
		for i := range accounts[1:] {
			sleepPartOfTotalTime(int64(len(accounts)), 6*time.Hour)
			startupTaskWorker(db, &accounts[1+i], &wg, retryTimes)
		}
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
	if info.VerNumber > lib_version {
		log.Println("Latest App version: ", info.VerNumber)
		log.Println("Need to upgrade!!!")
		return false
	}
	return true
}

func startupTaskWorker(db *database.DB, acc *accountSrv.Account, wg *sync.WaitGroup, retryTimes int) {
	log.Println("runAccountTask", acc.SchoolID, acc.StuNum)
	go func() {
		defer wg.Done()
		setAccountStatus(db, acc, accountSrv.StatusRunning)
		failCnt := 0
		forceUpdateSession := false
	execute:
		for failCnt < retryTimes {
			err := newTask(db, acc, forceUpdateSession).Exec()
			setAccountLastTime(db, acc, time.Now())
			switch err {
			case nil:
				log.Println("runAccountTask", acc.SchoolID, acc.StuNum, "has been completed Successfully.")
				setAccountLastResult(db, acc, accountSrv.RunSuccess)
				setAccountStatus(db, acc, accountSrv.StatusNormal)
			case ErrFinished:
				setAccountLastResult(db, acc, accountSrv.RunSuccess)
				setAccountStatus(db, acc, accountSrv.StatusFinished)
			case ssmt.ErrInvalidToken:
				forceUpdateSession = true
				failCnt++
				continue execute
			default:
				fmt.Println(acc.SchoolID, acc.StuNum, ": ", err.Error())
				setAccountLastResult(db, acc, accountSrv.RunErrorOccurred)
				setAccountStatus(db, acc, accountSrv.StatusSuspend)
			}
			break execute // exit normally
		}
	}()
}
