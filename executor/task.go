package main

import (
	"errors"
	"log"
	"math/rand"
	"time"

	ssmt "github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv/accLogSrv"
	"github.com/inkedawn/JKWXRunner-server/service/sessionSrv"
	"github.com/inkedawn/JKWXRunner-server/service/userCacheSrv"
	"github.com/inkedawn/JKWXRunner-server/viewFormat"
)

var (
	ErrFinished           = errors.New("已完成跑步，不需要再执行任务")
	ErrWrongLibVersion    = errors.New("错误的库版本")
	ErrCheatMarked        = errors.New("该帐号已被标记作弊！")
	ErrUnexpectedUserInfo = errors.New("帐号信息异常，可能是session存在问题。")
)

type task struct {
	dbSrv                         service.ICommonService
	acc                           *datamodels.Account
	enableRandomDistanceReduction bool
	forceUpdateSession            bool
}

func newTask(dbSrv service.ICommonService, acc *datamodels.Account, forceUpdateSession bool) *task {
	return &task{dbSrv: dbSrv, acc: acc, enableRandomDistanceReduction: true, forceUpdateSession: forceUpdateSession}
}

func (t *task) Exec() (err error) {
	log.Println("runAccountTask", t.acc.SchoolID, t.acc.StuNum)
	db := t.dbSrv.GetDB()
	acc := t.acc
	uid := t.acc.ID
	defer func() {
		if x := recover(); x != nil {
			err, _ = x.(error)
			if err == ssmt.ErrInvalidToken {
				accLogSrv.AddLogInfo(db, t.acc.ID, "Session失效.")
			}
		}
	}()
	var s *ssmt.Session
	if t.forceUpdateSession {
		s, err = sessionSrv.NewSession(db, *acc)
		if err != nil {
			accLogSrv.AddLogFail(db, uid, "更新Session失败："+dumpStruct(err))
			return err
		}
		accLogSrv.AddLogInfo(db, uid, "更新Session成功："+dumpStruct(s))
	} else {
		s, err = sessionSrv.SmartGetSession(db, *acc)
		if err != nil {
			accLogSrv.AddLogFail(db, uid, "创建Session失败："+dumpStruct(err))
			return err
		}
	}

	userInfo, err := userCacheSrv.GetCacheUserInfo(db, s.User.UserID)
	if err != nil {
		accLogSrv.AddLogFail(db, uid, "获取UserInfo失败："+dumpStruct(err))
		return err
	}
	if userInfo.UserRoleID == service.UserRole_Cheater {
		accLogSrv.AddLogInfo(db, uid, "检测到该帐号已被标记作弊！")
		// 从数据库取回的应当必定该字段有效
		if !acc.CheckCheatMarked.Valid {
			panic("标记作弊设定异常")
		}
		if acc.CheckCheatMarked.Bool {
			accLogSrv.AddLogFail(db, uid, "根据标记作弊设定。停止执行")
			return ErrCheatMarked
		}
	}
	if userInfo.Sex != "F" && userInfo.Sex != "M" {
		accLogSrv.AddLogFail(db, uid, "未知的性别：", userInfo.Sex)
		return ErrUnexpectedUserInfo
	}
	limit := ssmt.GetDefaultLimitParams(userInfo.Sex)

	r, err := recordResultBeforeRun(db, acc.ID, s)
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
	if info.VerNumber > libVersion {
		log.Println("Latest App version: ", info.VerNumber)
		log.Println("Need to upgrade!!!")
		return ErrWrongLibVersion
	}
	stillNeed := r.QualifiedDistance - r.ActualDistance
	if stillNeed < acc.RunDistance {
		// 接近完成，就不跑满
		limit.RandDistance.Min = stillNeed + 0.1
		limit.RandDistance.Max = stillNeed + 0.8
		accLogSrv.AddLogInfoF(db, uid, "即将完成。本次随机区间 %v~%v", viewFormat.DistanceFormat(limit.RandDistance.Min), viewFormat.DistanceFormat(limit.RandDistance.Max))
	} else if t.enableRandomDistanceReduction {
		// 一定几率不跑满，触发几率
		const (
			// the trigger rate is triggerRateN/triggerRateM
			triggerRateN = 2
			triggerRateM = 18
		)
		triggerRand := rand.Intn(triggerRateM)
		accLogSrv.AddLogDebug(db, uid, "triggerRand:", triggerRand, " triggerRateN:", triggerRateN, " triggerRateM:", triggerRateM)
		if triggerRand < triggerRateN {
			const (
				// the rate range is (0, maxMinusRate/reductionRateDivision)
				maxReductionRate      = 20
				reductionRateDivision = 100
			)
			reduceRate := float64(1+rand.Intn(maxReductionRate)) / reductionRateDivision
			limit.RandDistance.Max = limit.LimitSingleDistance.Min + (limit.RandDistance.Max-limit.LimitSingleDistance.Min)*(1-reduceRate)
			accLogSrv.AddLogInfoF(db, uid, "本次触发不跑满策略，比率%v，最终上限%v", reduceRate, viewFormat.DistanceFormat(limit.RandDistance.Max))
		}
	}
	records := ssmt.SmartCreateRecordsAfter(s.User.SchoolID, s.User.UserID, limit, acc.RunDistance, time.Now())
	err = uploadRecords(db, acc, s, records)
	if err != nil {
		// still record distance if upload records fail
		_, _ = recordResultAfterRun(db, acc.ID, s) // but if record fail also, let it go
		return err
	}
	// major operation has been completed successfully.
	r, err = recordResultAfterRun(db, acc.ID, s)
	if err != nil {
		accLogSrv.AddLogFailF(db, acc.ID, "结束后记录距离时遇到错误。", err)
		// only log but not return
		return nil
	}
	if shouldFinished(acc, r) {
		return ErrFinished
	}
	return nil
}
