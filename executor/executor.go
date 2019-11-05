package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
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

func startupTaskWorker(db *database.DB, acc *datamodels.Account, wg *sync.WaitGroup, retryTimes int) {
	log.Println("runAccountTask", acc.SchoolID, acc.StuNum)
	go func() {
		defer wg.Done()
		setAccountStatus(db, acc, service.AccountStatusRunning)
		failCnt := 0
		task := newTask(db, acc, false)
	execute:
		for failCnt < retryTimes {
			err := task.Exec()
			setAccountLastTime(db, acc, time.Now())
			switch err {
			case nil:
				log.Println("runAccountTask", acc.SchoolID, acc.StuNum, "has been completed Successfully.")
				setAccountLastResult(db, acc, service.TaskRunSuccess)
				setAccountStatus(db, acc, service.AccountStatusNormal) // 从Running状态恢复。
			case ssmt.ErrInvalidToken:
				task.forceUpdateSession = true
				failCnt++
				if failCnt < retryTimes {
					continue execute
				} else {
					setAccountLastResult(db, acc, service.TaskRunErrorOccurred)
					setAccountStatus(db, acc, service.AccountStatusSuspend)
				}
			case ErrFinished:
				setAccountLastResult(db, acc, service.TaskRunSuccess)
				setAccountStatus(db, acc, service.AccountStatusFinished) // 超出距离自动更改为结束状态。
			default:
				fmt.Println(acc.SchoolID, acc.StuNum, ": ", err.Error())
				setAccountLastResult(db, acc, service.TaskRunErrorOccurred)
				setAccountStatus(db, acc, service.AccountStatusSuspend) // 遇到未知错误将自动挂起
			}
			break execute // exit normally
		}
	}()
}
