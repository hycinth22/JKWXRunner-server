package main

import (
	"./model"
	sunshinemotion "./sunshinemotion"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
	"strconv"
)

func registerTicketRoute(router gin.IRouter) {
	router.GET("/allTicket", getAllTickets)
	router.POST("/ticket", newTicket)
	router.PUT("/ticket", updateTicket)
	router.DELETE("/ticket/:id", deleteTicket)

	router.GET("/ticket/log/:id", getTicketLog)
}

func getAllTickets(context *gin.Context) {
	list, err := model.GetAllTickets()
	if err != nil {
		context.Error(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, list)
}

func newTicket(context *gin.Context) {
	var ticket model.TicketWithSingleAccount
	if err := context.MustBindWith(&ticket, binding.JSON); err != nil {
		log.Println(context.ClientIP(), err.Error())
		return
	}

	session := sunshinemotion.CreateSession()
	if err := session.Login(ticket.Account.Username, "123", sunshinemotion.PasswordHash(ticket.Password)); err != nil {
		log.Println(context.ClientIP(), err.Error())
		return
	}
	ticket.Account.RemoteUserID = session.UserID
	ticket.Account.Profile = session.UserInfo

	id, err := model.AddTicket(&ticket)
	if err != nil {
		context.Error(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	model.SaveSession(ticket.Account.ID, session)

	context.JSON(http.StatusAccepted, struct{ ID uint }{ID: id})
}

func updateTicket(context *gin.Context) {
	var ticket model.TicketWithSingleAccount
	if err := context.MustBindWith(&ticket, binding.JSON); err != nil {
		log.Println(context.ClientIP(), err.Error())
		return
	}
	if err := ticket.Update(); err != nil {
		context.Error(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusAccepted)
}

func deleteTicket(context *gin.Context) {
	var err error
	idNum, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id := uint(idNum)
	err = model.DelTicketByID(id)
	if err != nil {
		context.Error(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusAccepted)
}

func getTicketLog(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// log.Println("request lookup log Ticket" , id)
	ticket, err := model.GetTicketByID(uint(id))
	if err != nil {
		context.Error(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// log.Println("ticket.Account.ID" , ticket.Account.ID, ticket.Account)
	list := model.GetLogs(ticket.Account.ID, 30)
	context.JSON(http.StatusOK, list)
}
