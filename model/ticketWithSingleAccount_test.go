package model

import (
	"strconv"
	"testing"
	"time"
)

func TestAddTicket(t *testing.T) {
	for i := 0; i < 5; i++ {
		identity := strconv.FormatInt(time.Now().Unix()+int64(i), 10)
		ticket := &TicketWithSingleAccount{
			Account: Account{
				RemoteUserID: time.Now().UnixNano() + int64(i),
				Username:     identity,
				Password:     "testPassword_for_" + identity,
				Distance:     0.111 + float64(i),
				RunResult: RunResult{
					LastTime:   time.Now(),
					LastStatus: StatusOK,
				},
			},
			Ticket: Ticket{
				OwnerID: 0,
				Contact: "testContact_for_" + identity,
				Memo:    "testMemo_for_" + identity,
			},
		}
		id, err := AddTicket(ticket, CachedUserInfo{})
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		t.Logf("Add id: %d", id)
	}
}

func TestListTickets(t *testing.T) {
	list, err := ListTickets()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	for _, ticket := range list {
		t.Logf("GetAll: %v\n", ticket)
	}
}

func TestGetTicketByID(t *testing.T) {
	identity := strconv.FormatInt(time.Now().UnixNano(), 10)
	id, err := AddTicket(&TicketWithSingleAccount{
		Account: Account{
			RemoteUserID: time.Now().UnixNano(),
			Username:     "ttttttttttttttttttttttttt" + identity,
			Password:     "testPassword_for_" + identity,
			Distance:     0.111,
			RunResult: RunResult{
				LastTime:   time.Now(),
				LastStatus: StatusOK,
			},
		},
		Ticket: Ticket{
			OwnerID: 0,
			Contact: "testContact_for_" + identity,
			Memo:    "testMemo_for_" + identity,
		}}, CachedUserInfo{})
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	ticket, err := GetTicketByID(id)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	if ticket.Password != "testPassword_for_"+identity {
		t.Fail()
		return
	}
	t.Logf("GetTicket: %v", ticket)
}

func TestTicketWithSingleAccount_Update(t *testing.T) {
	oldTicket, err := GetTicketByID(4)
	if err != nil {
		t.Fail()
		return
	}
	t.Log("old ", oldTicket.LastStatus)

	err = oldTicket.Save()

	newTicket := oldTicket
	newTicket.LastStatus = StatusCompleted
	t.Log("modify to ", oldTicket.LastStatus)
	newTicket.Save()

	resultTicket, err := GetTicketByID(4)
	if err != nil {
		t.Fail()
		return
	}
	t.Log("result", resultTicket.LastStatus)
	err = oldTicket.Save()
	if err != nil {
		t.Fatal("recover updated ticket failed", err)
		return
	}

	if newTicket.LastStatus != resultTicket.LastStatus {
		t.Fail()
		t.Log(newTicket)
		t.Log(resultTicket)
	}
}
