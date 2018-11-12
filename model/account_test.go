package model

import "testing"

func TestListAccountsTodayNotRun(t *testing.T) {
	a, err := ListAccountsTodayNotRun()
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
