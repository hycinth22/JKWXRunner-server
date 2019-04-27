package main

import (
	"errors"
	"fmt"
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/service/accountSrv"
	"github.com/inkedawn/JKWXFucker-server/service/accountSrv/accLogSrv"
	"github.com/inkedawn/JKWXFucker-server/service/sessionSrv"
	"github.com/inkedawn/JKWXFucker-server/service/userCacheSrv"
	"github.com/inkedawn/go-sunshinemotion"
	"log"
	"time"
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
		if acc.CheckCheatMarked && userInfo.UserRoleID == userCacheSrv.UserRole_Cheater {
			accLogSrv.AddLogFail(db, uid, "该帐号已被标记作弊！停止执行")
			return ErrCheatMarked
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
		if r.ActualDistance > r.QualifiedDistance {
			return ErrFinished
		}

		info, err := s.GetAppInfo()
		if err != nil {
			return err
		}
		if info.VerNumber != lib_version {
			log.Println("Latest App version: ", info.VerNumber)
			log.Println("Need to upgrade!!!")
			return ErrWrongLibVersion
		}

		records := ssmt.SmartCreateRecordsAfter(s.User.UserID, s.User.SchoolID, limit, acc.RunDistance, time.Now())
		err = uploadRecords(db, acc, s, records)

		if err != nil {
			_, _ = recordResultAfterRun(db, acc.ID, s) // if fail, let it go
			return err
		}
		r, err = recordResultAfterRun(db, acc.ID, s)
		if err != nil {
			return err
		}
		if r.ActualDistance > r.QualifiedDistance {
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
		_, err := s.GetRandRoute()
		if err != nil {
			accLogSrv.AddLogFail(db, uid, fmt.Sprintf("第%d条记录GetRandRoute失败：%#v。", n, err))
			return errors.New("GetRandRoute" + err.Error())
		}
		log.Println(n, r)
		log.Println("Sleep Util", r.EndTime)
		sleepUtil(r.EndTime)
		err = s.UploadRecord(r)
		if err != nil {
			accLogSrv.AddLogFail(db, uid, fmt.Sprintf("上传第%d条记录失败：%#v。 RecordDump: %s", n, err, dumpStructValue(r)))
			return err
		}
		accLogSrv.AddLogSuccess(db, uid, fmt.Sprintf("上传第%d条记录成功。 RecordDump: %s", n, dumpStructValue(r)))
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
