package main

import (
	"errors"
	"github.com/inkedawn/JKWXFucker-server/model"
	"github.com/inkedawn/JKWXFucker-server/view"
	sunshinemotion "github.com/inkedawn/go-sunshinemotion"
	"strconv"
	"time"
)

func RunForAccount(account *model.Account) model.RunResult {
	// (status model.Status, lastTime time.Time, lastDistance float64)
	s, err := model.GetSession(account.ID)
	if err != nil || s.UserExpirationTime.Before(time.Now()) {
		if err != nil {
			s = sunshinemotion.CreateSession()
			account.AddLog(time.Now(), model.LogTypeError, "获取session失败"+err.Error())
		}
		err := s.Login(account.Username, "123", sunshinemotion.PasswordHash(account.Password))
		if err != nil {
			account.AddLog(time.Now(), model.LogTypeError, "登录失败: "+err.Error())
			return model.RunResult{
				LastStatus:   model.StatusFail,
				LastTime:     time.Now(),
				LastDistance: 0.0,
			}
		}
		account.AddLog(time.Now(), model.LogTypeInfo, "登录成功")
		model.SaveSession(account.ID, s)
		randSleep(15*time.Second, 30*time.Second)
	}
	result, err := s.GetSportResult()
	if err == nil {
		account.AddLog(time.Now(), model.LogTypeInfo, "上传前已跑距离"+view.DistanceFormat(result.Distance))
		saveCachedUserInfo(account, model.CachedUserInfo{
			TotalDistance:     result.Distance,
			QualifiedDistance: result.Qualified,
		})
	} else {
		account.AddLog(time.Now(), model.LogTypeError, "上传前获取已跑信息失败")
	}
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
	records := sunshinemotion.SmartCreateRecords(account.RemoteUserID, s.LimitParams, account.Distance, time.Now())

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

func saveRunResult(account *model.Account, result model.RunResult) error {
	account.RunResult = result
	return model.UpdateAccount(account)
}

func saveCachedUserInfo(account *model.Account, info model.CachedUserInfo) error {
	account.CachedUserInfo = info
	return model.UpdateAccount(account)
}
