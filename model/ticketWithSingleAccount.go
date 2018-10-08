// 暂时只支持单account的Ticket
package model

import (
	"errors"
)

type TicketWithSingleAccount struct {
	Ticket
	Account
}


func GetAllTickets() (tickets []*TicketWithSingleAccount, err error) {
	simpleTickets := make([]Ticket, 0)
	accounts := make([]Account, 0)
	if err := db.Find(&simpleTickets).Related(&accounts, "TicketID").Error; err != nil {
		return nil, errors.New("GetAllTickets Failed:" + err.Error())
	}
	for _, ticket := range simpleTickets {
		// TODO: 优化
		var account Account
		if err := db.Model(&ticket).Related(&account, "TicketID").Error; err != nil {
			return nil, errors.New("GetAllTickets Failed:" + err.Error())
		}
		var one TicketWithSingleAccount
		one.Ticket = ticket
		one.Account = account
		tickets = append(tickets, &one)
	}
	return tickets, nil
}

func GetTicketByID(ticketID uint) (ticket *TicketWithSingleAccount, err error) {
	simpleTicket := Ticket{
		ID: ticketID,
	}
	var account Account
	if err := db.First(&simpleTicket).Related(&account, "TicketID").Error; err != nil {
		return nil, errors.New("GetTicketByID Failed:" + err.Error())
	}
	ticket = new(TicketWithSingleAccount)
	ticket.Ticket = simpleTicket
	ticket.Account = account
	return ticket, nil
}

// add to database
func AddTicket(ticket *TicketWithSingleAccount) (ticketID uint, err error) {
	if !db.NewRecord(&ticket.Ticket) {
		return 0, errors.New("parameter ticket.Ticket is not a new record")
	}
	if !db.NewRecord(&ticket.Account) {
		return 0, errors.New("parameter ticket.Account is not a new record")
	}
	tx := db.Begin()
	// create ticket
	if err := tx.Create(&ticket.Ticket).Error; err != nil {
		tx.Rollback()
		return 0, errors.New("Add ticket.Ticket Fail:" + err.Error())
	}
	// create account
	ticket.Account.TicketID = ticket.Ticket.ID
	if err := tx.Create(&ticket.Account).Error; err != nil {
		tx.Rollback()
		return 0, errors.New("Add ticket.Account Fail:" + err.Error())
	}
	tx.Commit()
	return ticket.Ticket.ID, nil
}

// del from database
func DelTicketByID(ticketID uint) (err error) {
	ticket := new(TicketWithSingleAccount)
	ticket.Ticket.ID = ticketID
	return ticket.Del()
}

// del from database
func (ticket *TicketWithSingleAccount) Del() (err error){
	if ticket.Ticket.ID == 0{
		return
	}
	tx := db.Begin()
	var account Account
	if err := tx.Model(&ticket).Related(&account).Error; err != nil {
		tx.Rollback()
		return errors.New("get Deleted ticket's account fail:" + err.Error())
	}
	if err := tx.Delete(&account).Error; err != nil {
		tx.Rollback()
		return errors.New("del account fail:" + err.Error())
	}
	if err := tx.Delete(&ticket).Error; err != nil {
		tx.Rollback()
		return errors.New("del ticket fail:" + err.Error())
	}
	tx.Commit()
	return nil
}

// write to database
func (ticket *TicketWithSingleAccount) Update()  (err error) {
	tx := db.Begin()
	if err := tx.Save(&ticket).Error; err != nil {
		tx.Rollback()
		return errors.New("update ticket fail:" + err.Error())
	}
	tx.Commit()
	return nil
}
