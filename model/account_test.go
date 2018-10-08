package model

import "testing"

func TestGetAllTodayNotRun(t *testing.T) {
	a, err := GetAllAccountsTodayNotRun()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	for _, account := range a {
		t.Log(account.ID, account.LastTime)
	}
}

