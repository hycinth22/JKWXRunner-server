package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv/accLogSrv"
	"github.com/inkedawn/JKWXRunner-server/viewFormat"
)

func uploadRecords(db *database.DB, sportResultSrv service.IUserSportResultService, acc *datamodels.Account, s *ssmt.Session, records []ssmt.Record) error {
	uid := acc.ID
	for i, r := range records {
		n := i + 1
		var recordNoX string
		if len(records) > 1 {
			recordNoX = fmt.Sprintf("第%d条记录", n)
		} else {
			recordNoX = "本次记录"
		}
		route, err := s.GetRandRoute()
		if err != nil {
			accLogSrv.AddLogFail(db, uid, fmt.Sprintf("%sGetRandRoute失败：%#v。", recordNoX, err))
			return errors.New("GetRandRoute" + err.Error())
		}
		accLogSrv.AddLogInfo(db, uid, fmt.Sprintf("%s获取到本次跑步路线：%+v。", recordNoX, route))
		log.Println(uid, n, r, "Sleep Util", r.EndTime)
		accLogSrv.AddLogInfo(db, uid, fmt.Sprintf("%s生成的记录是：%s", recordNoX, dumpStructValue(r)))
		accLogSrv.AddLogInfo(db, uid, fmt.Sprintf("%s需等待至%s。", recordNoX, viewFormat.TimeFormat(r.EndTime)))
		sleepUtil(r.EndTime)
		_, err = recordResultBeforeUpload(db, sportResultSrv, acc.ID, s)
		if err != nil {
			return err
		}
		err = s.UploadRecord(r)
		if err != nil {
			accLogSrv.AddLogFail(db, uid, fmt.Sprintf("上传%s失败：%#v", recordNoX, err))
			return err
		}
		accLogSrv.AddLogSuccess(db, uid, fmt.Sprintf("上传%s成功", recordNoX))
		_, err = recordResultAfterUpload(db, sportResultSrv, acc.ID, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func recordResultBeforeExec(db *database.DB, sportResultSrv service.IUserSportResultService, uid uint, s *ssmt.Session) (result *ssmt.SportResult, err error) {
	result, err = s.GetSportResult()
	if err != nil {
		accLogSrv.AddLogFail(db, uid, "执行前获取已跑信息失败："+err.Error())
		return nil, err
	}

	_ = sportResultSrv.SaveCacheSportResult(datamodels.CacheUserSportResultFromSSMTSportResult(*result, s.User.UserID, time.Now()))
	accLogSrv.AddLogInfo(db, uid, "执行前运动结果： "+dumpStructValue(*result))
	return
}

func recordResultAfterExec(db *database.DB, sportResultSrv service.IUserSportResultService, uid uint, s *ssmt.Session) (result *ssmt.SportResult, err error) {
	result, err = s.GetSportResult()
	if err != nil {
		accLogSrv.AddLogFail(db, uid, "执行后获取已跑信息失败："+err.Error())
		return nil, err
	}
	_ = sportResultSrv.SaveCacheSportResult(datamodels.CacheUserSportResultFromSSMTSportResult(*result, s.User.UserID, time.Now()))
	accLogSrv.AddLogInfo(db, uid, "执行后运动结果： "+dumpStructValue(*result))
	return
}

func recordResultBeforeUpload(db *database.DB, sportResultSrv service.IUserSportResultService, uid uint, s *ssmt.Session) (result *ssmt.SportResult, err error) {
	result, err = s.GetSportResult()
	if err != nil {
		accLogSrv.AddLogFail(db, uid, "上传前获取已跑信息失败："+err.Error())
		return nil, err
	}

	_ = sportResultSrv.SaveCacheSportResult(datamodels.CacheUserSportResultFromSSMTSportResult(*result, s.User.UserID, time.Now()))
	accLogSrv.AddLogInfo(db, uid, "上传前运动结果： "+dumpStructValue(*result))
	return
}

func recordResultAfterUpload(db *database.DB, sportResultSrv service.IUserSportResultService, uid uint, s *ssmt.Session) (result *ssmt.SportResult, err error) {
	result, err = s.GetSportResult()
	if err != nil {
		accLogSrv.AddLogFail(db, uid, "上传后获取已跑信息失败："+err.Error())
		return nil, err
	}
	_ = sportResultSrv.SaveCacheSportResult(datamodels.CacheUserSportResultFromSSMTSportResult(*result, s.User.UserID, time.Now()))
	accLogSrv.AddLogInfo(db, uid, "上传后运动结果： "+dumpStructValue(*result))
	return
}

func shouldFinished(acc *datamodels.Account, result *ssmt.SportResult) bool {
	if acc.FinishDistance == 0.0 {
		acc.FinishDistance = result.QualifiedDistance
	}
	return result.ActualDistance >= acc.FinishDistance
}

func setAccountStatus(db *database.DB, acc *datamodels.Account, status service.AccountStatus) {
	err := accountSrv.SetStatus(db, acc, status)
	if err != nil {
		log.Println("account ", acc.SchoolID, acc.StuNum, "failed to set status to", status)
	}
}

func setAccountLastTime(db *database.DB, acc *datamodels.Account, t time.Time) {
	err := accountSrv.SetLastTime(db, acc, t)
	if err != nil {
		log.Println("account ", acc.SchoolID, acc.StuNum, "failed to set lastTime to", t)
	}
}

func setAccountLastResult(db *database.DB, acc *datamodels.Account, r service.TaskRunResult) {
	err := accountSrv.SetLastResult(db, acc, r)
	if err != nil {
		log.Println("account ", acc.SchoolID, acc.StuNum, "failed to set lastResult to", r)
	}
}
