package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv/accLogSrv"
	"github.com/inkedawn/JKWXRunner-server/service/userCacheSrv"
	"github.com/inkedawn/JKWXRunner-server/viewFormat"
)

func uploadRecords(db *database.DB, acc *accountSrv.Account, s *ssmt.Session, records []ssmt.Record) error {
	uid := acc.ID
	for i, r := range records {
		n := i + 1
		var recordNoPlaceHolder string
		if len(records) > 1 {
			recordNoPlaceHolder = fmt.Sprintf("第%d条记录", n)
		} else {
			recordNoPlaceHolder = "本次记录"
		}
		route, err := s.GetRandRoute()
		if err != nil {
			accLogSrv.AddLogFail(db, uid, fmt.Sprintf("%sGetRandRoute失败：%#v。", recordNoPlaceHolder, err))
			return errors.New("GetRandRoute" + err.Error())
		}
		accLogSrv.AddLogInfo(db, uid, fmt.Sprintf("%s获取到本次跑步路线：%+v。", recordNoPlaceHolder, route))
		log.Println(uid, n, r, "Sleep Util", r.EndTime)
		accLogSrv.AddLogInfo(db, uid, fmt.Sprintf("%s生成的记录是：%s", recordNoPlaceHolder, dumpStructValue(r)))
		accLogSrv.AddLogInfo(db, uid, fmt.Sprintf("%s需等待至%s。", recordNoPlaceHolder, viewFormat.TimeFormat(r.EndTime)))
		sleepUtil(r.EndTime)
		err = s.UploadRecord(r)
		if err != nil {
			accLogSrv.AddLogFail(db, uid, fmt.Sprintf("上传%s失败：%#v", recordNoPlaceHolder, err))
			return err
		}
		accLogSrv.AddLogSuccess(db, uid, fmt.Sprintf("上传%s成功", recordNoPlaceHolder))
	}
	return nil
}

func recordResultBeforeRun(db *database.DB, uid uint, s *ssmt.Session) (result *ssmt.SportResult, err error) {
	result, err = s.GetSportResult()
	if err != nil {
		accLogSrv.AddLogFail(db, uid, "上传前获取已跑信息失败："+err.Error())
		return nil, err
	}
	_ = userCacheSrv.SaveCacheSportResult(db, userCacheSrv.FromSSMTSportResult(*result, s.User.UserID, time.Now()))
	accLogSrv.AddLogInfo(db, uid, "上传前运动结果： "+dumpStructValue(*result))
	return
}

func recordResultAfterRun(db *database.DB, uid uint, s *ssmt.Session) (result *ssmt.SportResult, err error) {
	result, err = s.GetSportResult()
	if err != nil {
		accLogSrv.AddLogFail(db, uid, "上传后获取已跑信息失败："+err.Error())
		return nil, err
	}
	_ = userCacheSrv.SaveCacheSportResult(db, userCacheSrv.FromSSMTSportResult(*result, s.User.UserID, time.Now()))
	accLogSrv.AddLogInfo(db, uid, "上传后运动结果： "+dumpStructValue(*result))
	return
}

func shouldFinished(acc *accountSrv.Account, result *ssmt.SportResult) bool {
	if acc.FinishDistance == 0.0 {
		acc.FinishDistance = result.QualifiedDistance
	}
	return result.ActualDistance >= acc.FinishDistance
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
