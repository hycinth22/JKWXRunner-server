// 提供Account的log记录
package accLogSrv

import (
	"fmt"
	"os"
	"time"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/datamodels"
)

type Log = datamodels.AccountLog

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

func ListLogsForUID(db database.TX, uid uint, offset, limit uint) (logs []datamodels.AccountLog, err error) {
	err = db.Where("uid = ?", uid).Offset(offset).Limit(limit).Order("time desc").Find(&logs).Error
	return
}

func CountLogsForUID(db database.TX, uid uint) (n int, err error) {
	err = db.Model(&Log{}).Where("uid = ?", uid).Count(&n).Error
	return
}

func AddLogNow(db database.TX, uid uint, logType Type, text string) {
	AddLog(db, uid, time.Now(), logType, text)
}

func AddLog(db database.TX, uid uint, time time.Time, logType Type, text string) {
	err := db.Create(&Log{UID: uid, Time: time, Type: logType, Content: text}).Error
	if err != nil {
		reportErr(fmt.Sprintf("[%s] UID: %d, Time:%v, logType:%s, Text:%s", serviceName, uid, time, logType, text))
	}
}

/*
	the following functions are handy-function,
	reducing the parameter number for convenience..
*/

func AddLogSuccess(db database.TX, uid uint, values ...interface{}) {
	AddLogNow(db, uid, TypeSuccess, fmt.Sprint(values...))
}
func AddLogFail(db database.TX, uid uint, values ...interface{}) {
	AddLogNow(db, uid, TypeFail, fmt.Sprint(values...))
}
func AddLogInfo(db database.TX, uid uint, values ...interface{}) {
	AddLogNow(db, uid, TypeInfo, fmt.Sprint(values...))
}
func AddLogDebug(db database.TX, uid uint, values ...interface{}) {
	AddLogNow(db, uid, TypeDebug, fmt.Sprint(values...))
}

func AddLogSuccessF(db database.TX, uid uint, format string, values ...interface{}) {
	AddLogNow(db, uid, TypeSuccess, fmt.Sprintf(format, values...))
}
func AddLogFailF(db database.TX, uid uint, format string, values ...interface{}) {
	AddLogNow(db, uid, TypeFail, fmt.Sprintf(format, values...))
}
func AddLogInfoF(db database.TX, uid uint, format string, values ...interface{}) {
	AddLogNow(db, uid, TypeInfo, fmt.Sprintf(format, values...))
}
func AddLogDebugF(db database.TX, uid uint, format string, values ...interface{}) {
	AddLogNow(db, uid, TypeDebug, fmt.Sprintf(format, values...))
}
