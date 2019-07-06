// 提供Account的log记录
package accLogSrv

import (
	"fmt"
	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/database/model"
	"os"

	"time"
)

type Log = model.AccountLog

type Type = string

const (
	TypeSuccess Type = "success"
	TypeFail    Type = "fail"
	TypeInfo    Type = "info"
	TypeDebug   Type = "debug"
)

// if any error occurred, write into ErrorOutput
// if writing into ErrorOutput fails, println
var ErrorOutput = os.Stderr

const serviceName = "accLogSrv"

func reportErr(msg string) {
	if _, err := ErrorOutput.WriteString(msg); err != nil {
		println("writing into ErrorOutput fails, try println")
		println(msg)
	}
}

func ListLogsForUID(db *database.DB, uid uint, offset, limit uint) (logs []model.AccountLog, err error) {
	err = db.Where("uid = ?", uid).Offset(offset).Limit(limit).Order("time desc").Find(&logs).Error
	return
}

func CountLogsForUID(db *database.DB, uid uint) (n int, err error) {
	err = db.Model(&Log{}).Where("uid = ?", uid).Count(&n).Error
	return
}

func AddLogNow(db *database.DB, uid uint, logType Type, text string) {
	AddLog(db, uid, time.Now(), logType, text)
}

func AddLog(db *database.DB, uid uint, time time.Time, logType Type, text string) {
	err := db.Create(&Log{UID: uid, Time: time, Type: logType, Content: text}).Error
	if err != nil {
		reportErr(fmt.Sprintf("[%s] UID: %d, Time:%v, logType:%s, Text:%s", serviceName, uid, time, logType, text))
	}
}

/*
	the following functions are handy-function,
	reducing the parameter number for convenience..
*/

func AddLogSuccess(db *database.DB, uid uint, text string) {
	AddLogNow(db, uid, TypeSuccess, text)
}
func AddLogFail(db *database.DB, uid uint, text string) {
	AddLogNow(db, uid, TypeFail, text)
}
func AddLogInfo(db *database.DB, uid uint, text string) {
	AddLogNow(db, uid, TypeInfo, text)
}
func AddLogDebug(db *database.DB, uid uint, text string) {
	AddLogNow(db, uid, TypeDebug, text)
}
