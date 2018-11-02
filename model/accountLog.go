package model

import (
	"errors"
	"log"
	"time"
)

type AccountLog struct {
	AccountID uint      `gorm:"INDEX:account_id" json:"accountID"`
	Time      time.Time `json:"time"`
	Type      LogType   `json:"type"`
	Content   string    `json:"content"`
}

type LogType uint

const (
	LogTypeSuccess = iota
	LogTypeError
	LogTypeInfo
)

func addLog(log *AccountLog) (err error) {
	if err := db.Create(log).Error; err != nil {
		return errors.New("addLog fail: " + err.Error())
	}
	return nil
}

// the latest record of the previous n at offset n
func GetLogs(AccountID uint, offset, n int) (logs []AccountLog) {
	// log.Println("lookup log where id", AccountID)
	if err := db.Model(&AccountLog{}).Where("account_id = ?", AccountID).Order("Time DESC").Offset(offset).Limit(n).Find(&logs).Error; err != nil && !db.RecordNotFound() {
		log.Println("GetLogs", err.Error())
		return nil
	}
	return logs
}

// the latest record of the previous n
func GetLogsTotalNum(AccountID uint) uint {
	var count uint
	if err := db.Model(&AccountLog{}).Where("account_id = ?", AccountID).Count(&count).Error; err != nil && !db.RecordNotFound() {
		log.Println("GetLogsTotalNum Error", err.Error())
		return 0
	}
	return count
}
