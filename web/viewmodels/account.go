package viewmodels

import (
	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/viewFormat"
)

type Account struct {
	ID               uint
	OwnerID          int
	CreatedAt        string
	SchoolID         int64
	StuNum           string
	Memo             string
	Status           string
	RunDistance      float64
	StartDistance    float64
	FinishDistance   float64
	CurrentDistance  float64
	CheckCheatMarked bool
	LastResult       string
	LastTime         string
}

func NewAccount(acc *datamodels.Account, currentDistance float64) *Account {
	return &Account{
		ID:               acc.ID,
		OwnerID:          acc.OwnerID,
		CreatedAt:        viewFormat.TimeFormat(acc.CreatedAt),
		SchoolID:         acc.SchoolID,
		StuNum:           acc.StuNum,
		Memo:             acc.Memo,
		Status:           acc.Status,
		RunDistance:      acc.RunDistance,
		StartDistance:    acc.StartDistance,
		FinishDistance:   acc.FinishDistance,
		CurrentDistance:  currentDistance,
		CheckCheatMarked: acc.CheckCheatMarked.Valid && acc.CheckCheatMarked.Bool,
		LastResult:       acc.LastResult,
		LastTime:         viewFormat.TimeFormat(acc.LastTime),
	}
}
