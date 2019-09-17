package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service/userCacheSrv"
)

func QueryCacheUserSportResult(ctx *gin.Context) {
	inputUID := ctx.Param("uid")
	if inputUID == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	uid, err := strconv.ParseUint(inputUID, 10, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db := database.GetDB()
	result, err := userCacheSrv.GetLocalUserCacheSportResult(db, uint(uid))
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}
