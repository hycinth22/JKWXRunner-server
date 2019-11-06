package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	ssmt "github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
)

type worker struct {
	dbSrv      service.ICommonService
	tasks      []*task
	retryTimes int
	wg         *sync.WaitGroup
	resolved   bool
}

func (w *worker) work() {
	for i, task := range w.tasks {
		if i != 0 {
			sleepPartOfTotalTime(int64(len(w.tasks)), 6*time.Hour)
		}
		w.ExecTask(task)
	}
}

func (w *worker) ExecTask(task *task) {
	db := w.dbSrv.GetDB()
	acc := task.acc
	setAccountStatus(db, acc, service.AccountStatusRunning)
	failCnt := 0
execute:
	for failCnt < w.retryTimes {
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
			if failCnt < w.retryTimes {
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
}

func startupTaskWorker(dbSrv service.ICommonService, accounts []*datamodels.Account, wg *sync.WaitGroup, retryTimes int) *worker {
	w := &worker{dbSrv: dbSrv, tasks: nil, retryTimes: retryTimes, wg: wg}
	for _, acc := range accounts {
		w.tasks = append(w.tasks, newTask(dbSrv, acc, false))
	}
	go func() {
		defer wg.Done()
		w.work()
	}()
	return w
}
