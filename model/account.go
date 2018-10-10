package model

import (
	"time"
)

type Account struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	TicketID  uint       `gorm:"UNIQUE;NOT NULL" json:"-"`

	RemoteUserID int64   `gorm:"UNIQUE;NOT NULL" json:"userID"`
	Username     string  `gorm:"UNIQUE;NOT NULL;UNIQUE_INDEX:username" json:"username"`
	Password     string  `json:"password"`
	Distance     float64 `json:"distance"`

	LastStatus   Status    `json:"lastStatus"`
	LastDistance float64   `json:"lastDistance"`
	LastTime     time.Time `json:"lastTime"`
}
type Status uint

const (
	StatusOK Status = iota
	StatusPartialFail
	StatusFail
)

func GetAllAccountsTodayNotRun() (accounts []*Account, err error) {
	accounts = make([]*Account, 0)
	now := time.Now()
	err = db.Where("last_time < ?", time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)).Find(&accounts).Error
	return
}

func GetAllAccounts() (accounts []*Account, err error) {
	accounts = make([]*Account, 0)
	err = db.Find(&accounts).Error
	return
}

func GetAccountByUsername(username string) (account *Account, err error) {
	account = new(Account)
	err = db.Where("username = ?", username).First(account).Error
	return
}

func UpdateAccount(account *Account) (err error) {
	db.Save(account)
	return db.Error
}

func (acc *Account) GetLogs(n int) (logs []AccountLog) {
	return GetLogs(acc.ID, n)
}

func (acc *Account) AddLog(time time.Time, logType LogType, content string) (err error) {
	return addLog(&AccountLog{
		AccountID: acc.ID,
		Time:      time,
		Type:      logType,
		Content:   content,
	})
}
