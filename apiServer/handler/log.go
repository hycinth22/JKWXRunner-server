package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/service/accountSrv/accLogSrv"
	"net/http"
	"strconv"
)

func ListLogsByUID(context *gin.Context) {
	uid, err := strconv.ParseUint(context.Param("uid"), 10, 64)
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
	context.JSON(http.StatusOK, logs)
}
