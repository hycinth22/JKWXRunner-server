package main

import (
	"../model"
	sunshinemotion "../sunshinemotion"
	"../view"
	"errors"
	"strconv"
	"time"
)

func RunForAccount(account *model.Account) RunResult {
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
			return RunResult{
				status:       model.StatusFail,
				lastTime:     time.Now(),
				lastDistance: 0.0,
			}
		}
		account.AddLog(time.Now(), model.LogTypeInfo, "登录成功: "+err.Error())
		model.SaveSession(account.ID, s)
		randSleep(15*time.Second, 30*time.Second)
	}
	failCnt := 0
	lastDistance := 0.0
	lastTime := time.Now()
	records := sunshinemotion.SmartCreateRecords(account.RemoteUserID, s.LimitParams, account.Distance, time.Now())
	for i, record := range records {
		result, err := s.GetSportResult()
		if err == nil {
			account.AddLog(time.Now(), model.LogTypeInfo, "第"+strconv.Itoa(i+1)+"条记录上传前已跑距离"+view.DistanceFormat(result.Distance))
		} else {
			account.AddLog(time.Now(), model.LogTypeError, "第"+strconv.Itoa(i+1)+"条记录上传前获取已跑信息失败")
		}
		if !Debug {
			err = s.UploadRecord(record)
		} else {
			err = errors.New("test Error")
		}
		if err != nil {
			failCnt++
			account.AddLog(time.Now(), model.LogTypeError, "第"+strconv.Itoa(i+1)+"条记录上传失败，原因是："+err.Error())
		} else {
			lastDistance += record.Distance
			lastTime = record.EndTime
			account.AddLog(time.Now(), model.LogTypeSuccess, "第"+strconv.Itoa(i+1)+"条记录上传成功")
		}
		result, err = s.GetSportResult()
		if err == nil {
			account.AddLog(time.Now(), model.LogTypeInfo, "第"+strconv.Itoa(i+1)+"条记录上传后已跑距离"+view.DistanceFormat(result.Distance))
		} else {
			account.AddLog(time.Now(), model.LogTypeError, "第"+strconv.Itoa(i+1)+"条记录上传后获取已跑信息失败")
		}
	}
	var status model.Status
	if failCnt == 0 {
		status = model.StatusOK
	} else if failCnt < len(records) {
		status = model.StatusPartialFail
	} else {
		status = model.StatusFail
	}
	return RunResult{
		status:       status,
		lastTime:     lastTime,
		lastDistance: lastDistance,
	}
}

type RunResult struct {
	status       model.Status
	lastTime     time.Time
	lastDistance float64
}
