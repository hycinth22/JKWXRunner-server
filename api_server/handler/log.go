package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv/accLogSrv"
)

func ListLogsByID(context *gin.Context) {
	uid, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		_ = context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	offset, err := strconv.ParseUint(context.Query("offset"), 10, 64)
	if err != nil {
		_ = context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	limit, err := strconv.ParseUint(context.Query("limit"), 10, 64)
	if err != nil {
		_ = context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	logs, err := accLogSrv.ListLogsForUID(database.GetDB(), uint(uid), uint(offset), uint(limit))
	if err != nil {
		_ = context.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	totalAmount, err := accLogSrv.CountLogsForUID(database.GetDB(), uint(uid))
	if err != nil {
		_ = context.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, struct {
		TotalAmount int
		Items       []accLogSrv.Log
	}{totalAmount, logs})
}
