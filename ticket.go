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
	router.GET("/allTicket", listTickets)
	router.POST("/ticket", newTicket)
	router.PUT("/ticket", updateTicket)
	router.DELETE("/ticket/:id", deleteTicket)

	router.GET("/ticket/:id/log", listTicketLogs)
}

func listTickets(context *gin.Context) {
	list, err := model.ListTickets()
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

	result, err := session.GetSportResult()
	if err != nil {
		context.Error(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	id, err := model.AddTicket(&ticket, model.CachedUserInfo{
		TotalDistance:     result.Distance,
		QualifiedDistance: result.Qualified,
	})
	if err != nil {
		context.Error(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	model.SaveSession(ticket.Account.ID, session)

	context.JSON(http.StatusCreated, struct{ ID uint }{ID: id})
}

func updateTicket(context *gin.Context) {
	var ticket model.TicketWithSingleAccount
	if err := context.MustBindWith(&ticket, binding.JSON); err != nil {
		log.Println(context.ClientIP(), err.Error())
		return
	}
	if err := ticket.Save(); err != nil {
		context.Error(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, ticket)
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
	context.Status(http.StatusOK)
}

func listTicketLogs(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	offsetParam := context.Query("offset")
	numParam := context.Query("num")
	var (
		offset int
		num    int
	)
	if offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil || offset < 0 {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}
	} else {
		offset = 0
	}
	if numParam != "" {
		num, err = strconv.Atoi(numParam)
		if err != nil || num < 0 {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if num > 100 {
			num = 100
		}
	} else {
		num = 10
	}

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
	list := model.ListLogsForAccount(ticket.Account.ID, offset, num)
	total := model.CountLogsForAccount(ticket.Account.ID)
	context.JSON(http.StatusOK, struct {
		Total uint               `json:"total"`
		Logs  []model.AccountLog `json:"logs"`
	}{total, list})
}
