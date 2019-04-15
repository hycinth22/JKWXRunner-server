package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/service/userCacheSrv"
	"net/http"
	"strconv"
)

func QueryCacheUserSportResult(ctx *gin.Context) {
	inputUserID := ctx.Query("remoteUserID")
	if inputUserID == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	userID, err := strconv.ParseInt(inputUserID, 10, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	result, err := userCacheSrv.GetCacheSportResult(database.GetDB(), userID)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}
