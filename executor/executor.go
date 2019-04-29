package main

import (
	"fmt"
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/service/accountSrv"
	"github.com/inkedawn/go-sunshinemotion"
	"log"
	"os"
	"sync"
	"time"
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

	var wg sync.WaitGroup
	wg.Add(len(accounts))
	if len(accounts) >= 1 {
		startupTaskWorker(db, &accounts[0], &wg)
		for i := range accounts[1:] {
			sleepPartOfTotalTime(int64(len(accounts)), 6*time.Hour)
			startupTaskWorker(db, &accounts[1+i], &wg)
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

func startupTaskWorker(db *database.DB, acc *accountSrv.Account, wg *sync.WaitGroup) {
	log.Println("runAccountTask", acc.SchoolID, acc.StuNum)
	go func() {
		defer wg.Done()
		executeTask(db, acc, 3)
	}()
}

func executeTask(db *database.DB, acc *accountSrv.Account, retryTimes uint) {
	if retryTimes < 0 {
		return
	}
	err := runAccountTask(db, acc)
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
		executeTask(db, acc, retryTimes-1)
	default:
		fmt.Println(acc.SchoolID, acc.StuNum, ": ", err.Error())
		setAccountLastResult(db, acc, accountSrv.RunErrorOccurred)
		setAccountStatus(db, acc, accountSrv.StatusSuspend)
	}
}

func setAccountStatus(db *database.DB, acc *accountSrv.Account, status accountSrv.Status) {
	err := accountSrv.SetStatus(db, acc, status)
	if err != nil {
		log.Println("account ", acc.SchoolID, acc.StuNum, "failed to set status to", status)
	}
}

func setAccountLastTime(db *database.DB, acc *accountSrv.Account, t time.Time) {
	err := accountSrv.SetLastTime(db, acc, t)
	if err != nil {
		log.Println("account ", acc.SchoolID, acc.StuNum, "failed to set lastTime to", t)
	}
}

func setAccountLastResult(db *database.DB, acc *accountSrv.Account, r accountSrv.RunResult) {
	err := accountSrv.SetLastResult(db, acc, r)
	if err != nil {
		log.Println("account ", acc.SchoolID, acc.StuNum, "failed to set lastResult to", r)
	}
}
