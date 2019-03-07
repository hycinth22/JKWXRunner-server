package main

import (
	"errors"
	"github.com/inkedawn/JKWXFucker-server/model"
	"github.com/inkedawn/JKWXFucker-server/view"
	sunshinemotion "github.com/inkedawn/go-sunshinemotion"
	"log"
	"strconv"
	"time"
	"fmt"
)

func LoginForAccount(account *model.Account) (*sunshinemotion.Session, error) {
	// first, fetch session from store
	s, err := model.GetSession(account.ID)
	if err == nil && s.UserExpirationTime.After(time.Now()) {
		account.AddLog(time.Now(), model.LogTypeInfo, "Session" + fmt.Sprintf("%v", s))
		_, err := s.GetSportResult()
		if err == sunshinemotion.ErrTokenExpired {
			account.AddLog(time.Now(), model.LogTypeInfo, "Token过期，尝试重新登录")
		}
		if err == nil {
			return s, nil
		}
	} else if err != nil {
		s = sunshinemotion.CreateSession()
		account.AddLog(time.Now(), model.LogTypeError, "获取session失败"+err.Error())
	}

	// not store or expired, to call login
	err = s.LoginEx(account.Username, "123", sunshinemotion.PasswordHash(account.Password), account.RemoteSchoolID)
	if err != nil {
		account.AddLog(time.Now(), model.LogTypeError, "登录失败: "+err.Error())
		return s, errors.New("登录失败")
	} else {
		account.AddLog(time.Now(), model.LogTypeInfo, "登录成功。" + "Session Dump：" + fmt.Sprintf("%v", s))
		log.Println("LoginForAccount, user", s.UserInfo.StudentNumber, s)
	}

	// save it to store
	model.SaveSession(account.ID, s)
	return s, nil
}

func RunForAccount(account *model.Account) model.RunResult {
	// (status model.Status, lastTime time.Time, lastDistance float64)
	s, err := LoginForAccount(account)
	if err != nil {
		return model.RunResult{
			LastStatus:   model.StatusFail,
			LastTime:     time.Now(),
			LastDistance: 0.0,
		}
	}

	result, err := s.GetSportResult()
	if err != nil {
		account.AddLog(time.Now(), model.LogTypeError, "上传前获取已跑信息失败" + err.Error())
		return model.RunResult{
			LastStatus:   model.StatusFail,
			LastTime:     time.Now(),
			LastDistance: 0.0,
		}
	}
	account.AddLog(time.Now(), model.LogTypeInfo, "上传前已跑距离"+view.DistanceFormat(result.Distance))
	saveCachedUserInfo(account, model.CachedUserInfo{
		TotalDistance:     result.Distance,
		QualifiedDistance: result.Qualified,
	})
	if result.Distance > result.Qualified {
		return model.RunResult{
			LastStatus:   model.StatusCompleted,
			LastTime:     time.Now(),
			LastDistance: 0.0,
		}
	}
	failCnt := 0
	lastDistance := 0.0
	lastTime := time.Now()
	limit := getLimitParamsForSmartCreateRecords(s, result.Distance, result.Qualified)
	records := sunshinemotion.SmartCreateRecords(account.RemoteUserID, account.RemoteSchoolID, limit, account.Distance, time.Now())

	for i, record := range records {
		if !Debug {
			err = s.UploadRecord(record)
		} else {
			err = errors.New("test Error")
		}
		if err != nil {
			failCnt++
			account.AddLog(time.Now(), model.LogTypeError, "第"+strconv.Itoa(i+1)+"条记录上传失败，原因是：\n"+err.Error())
		} else {
			lastDistance += record.Distance
			lastTime = record.EndTime
			account.AddLog(time.Now(), model.LogTypeSuccess, "第"+strconv.Itoa(i+1)+"条记录上传成功，"+
				"距离"+view.DistanceFormat(record.Distance)+"公里。\n"+
				"起始时间"+view.TimeFormat(record.BeginTime)+"\n"+
				"结束时间"+view.TimeFormat(record.EndTime))
		}

	}

	result, err = s.GetSportResult()
	if err == nil {
		account.AddLog(time.Now(), model.LogTypeInfo, "上传后已跑距离"+view.DistanceFormat(result.Distance))
		saveCachedUserInfo(account, model.CachedUserInfo{
			TotalDistance:     result.Distance,
			QualifiedDistance: result.Qualified,
		})
	} else {
		account.AddLog(time.Now(), model.LogTypeError, "上传后获取已跑信息失败")
	}

	result, err = s.GetSportResult()
	var status model.Status
	if err == nil && result.Distance > result.Qualified {
		account.AddLog(time.Now(), model.LogTypeInfo, "已达标")
		status = model.StatusCompleted
	} else if failCnt == 0 {
		status = model.StatusOK
	} else if failCnt < len(records) {
		status = model.StatusPartialFail
	} else {
		status = model.StatusFail
	}
	return model.RunResult{
		LastStatus:   status,
		LastTime:     lastTime,
		LastDistance: lastDistance,
	}
}

func getLimitParamsForSmartCreateRecords(s *sunshinemotion.Session, hasDistance float64, qualifiedDistance float64) *sunshinemotion.LimitParams {
	limit := s.LimitParams
	stillNeed := qualifiedDistance - hasDistance
	if stillNeed > s.LimitParams.LimitTotalDistance.Min && stillNeed < s.LimitParams.LimitTotalDistance.Max {
		limit.LimitTotalDistance.Max = stillNeed + 0.1
	}
	return limit
}

func saveRunResult(account *model.Account, result model.RunResult) error {
	account.RunResult = result
	return model.UpdateAccount(account)
}

func saveCachedUserInfo(account *model.Account, info model.CachedUserInfo) error {
	account.CachedUserInfo = info
	return model.UpdateAccount(account)
}
