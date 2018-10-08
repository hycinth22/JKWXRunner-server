package model

import (
	"errors"
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

// the latest record of the previous n
func GetLogs(AccountID uint, n int) (logs []AccountLog) {
	// log.Println("lookup log where id", AccountID)
	if err := db.Model(&Account{
		ID: AccountID,
	}).Where("account_id = ?", AccountID).Order("Time DESC").Limit(n).Find(&logs).Error; err != nil {
		// log.Println(err.Error())
		return nil
	}
	return
}
