package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv"
)

func UpdateAccountStatus(ctx *gin.Context) {
	var (
		payload struct {
			Status string
		}
		err error
	)
	db := database.GetDB()
	tid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	id := uint(tid)
	err = ctx.Bind(&payload)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	acc, err := accountSrv.GetAccount(db, id)
	if err == accountSrv.ErrNoAccount {
		ctx.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	acc.Status = payload.Status
	err = accountSrv.SaveAccount(db, acc)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusAccepted)
}
