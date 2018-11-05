package model

import "testing"

func TestGetAllTodayNotRun(t *testing.T) {
	a, err := GetAllAccountsTodayNotRun()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	for _, account := range a {
		if account.LastStatus >= ExecStatusEndDelim {
			t.Fail()
		}
		t.Log(account.ID, account.LastTime)
	}
}
