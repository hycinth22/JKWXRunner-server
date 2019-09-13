package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv/accLogSrv"
	"github.com/inkedawn/JKWXRunner-server/service/sessionSrv"
	"github.com/inkedawn/JKWXRunner-server/service/userCacheSrv"
	"github.com/inkedawn/JKWXRunner-server/viewFormat"
)

var (
	ErrFinished        = errors.New("已完成跑步，不需要再执行任务")
	ErrWrongLibVersion = errors.New("错误的库版本")
	ErrCheatMarked     = errors.New("该帐号已被标记作弊！")
)

func runAccountTask(db *database.DB, acc *accountSrv.Account) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err, _ = x.(error)
		}
	}()
	failCnt := 0
execute:
	for failCnt < 2 {
		uid := acc.ID
		s, err := sessionSrv.SmartGetSession(db, *acc)
		if err != nil {
			accLogSrv.AddLogFail(db, uid, "创建Session失败："+dumpStruct(err))
			return err
		}

		userInfo, err := userCacheSrv.GetCacheUserInfo(db, s.User.UserID)
		if err != nil {
			accLogSrv.AddLogFail(db, uid, "获取UserInfo失败："+dumpStruct(err))
			return err
		}
		if userInfo.UserRoleID == userCacheSrv.UserRole_Cheater {
			accLogSrv.AddLogInfo(db, uid, "检测到该帐号已被标记作弊！")
			if acc.CheckCheatMarked {
				accLogSrv.AddLogFail(db, uid, "根据标记作弊设定。停止执行")
				return ErrCheatMarked
			}
		}
		limit := ssmt.GetDefaultLimitParams(userInfo.Sex)

		r, err := recordResultBeforeRun(db, acc.ID, s)
		if err == ssmt.ErrInvalidToken {
			accLogSrv.AddLogInfo(db, uid, "Session失效，尝试更新Session。Old Session Dump: %s"+dumpStruct(*s))
			err = sessionSrv.UpdateSession(db, *acc)
			if err != nil {
				accLogSrv.AddLogFail(db, uid, "更新Session失败："+dumpStruct(err))
				return err
			}
			// Retry
			failCnt++
			continue execute
		}
		if err != nil {
			return err
		}
		if shouldFinished(acc, r) {
			return ErrFinished
		}

		info, err := s.GetAppInfo()
		if err != nil {
			return err
		}
		if info.VerNumber > lib_version {
			log.Println("Latest App version: ", info.VerNumber)
			log.Println("Need to upgrade!!!")
			return ErrWrongLibVersion
		}
		stillNeed := r.QualifiedDistance - r.ActualDistance
		if stillNeed < limit.LimitSingleDistance.Max {
			limit.RandDistance.Min = stillNeed + 0.1
			limit.RandDistance.Max = stillNeed + 0.8
		}
		records := ssmt.SmartCreateRecordsAfter(s.User.SchoolID, s.User.UserID, limit, acc.RunDistance, time.Now())
		err = uploadRecords(db, acc, s, records)

		if err != nil {
			_, _ = recordResultAfterRun(db, acc.ID, s) // if fail, let it go
			return err
		}
		r, err = recordResultAfterRun(db, acc.ID, s)
		if err != nil {
			return err
		}
		if shouldFinished(acc, r) {
			return ErrFinished
		}
		break execute
	}
	return nil
}

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
