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
				LastTime:     time.Now(),
				LastStatus:   StatusOK,
			},
			Ticket: Ticket{
				OwnerID: 0,
				Contact: "testContact_for_" + identity,
				Memo:    "testMemo_for_" + identity,
			},
		}
		id, err := AddTicket(ticket)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		t.Logf("Add id: %d", id)
	}
}

func TestGetAll(t *testing.T) {
	list, err := GetAllTickets()
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
			LastTime:     time.Now(),
			LastStatus:   StatusOK,
		},
		Ticket: Ticket{
			OwnerID: 0,
			Contact: "testContact_for_" + identity,
			Memo:    "testMemo_for_" + identity,
		}},
	)
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
	if ticket.Password !=  "testPassword_for_" + identity{
		t.Fail()
		return
	}
	t.Logf("GetTicket: %v", ticket)
}
