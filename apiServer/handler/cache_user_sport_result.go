package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/service/userCacheSrv"
	"github.com/inkedawn/JKWXFucker-server/service/userIDRelationSrv"
	"net/http"
	"strconv"
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
	remoteUID, err := userIDRelationSrv.GetRemoteUserID(db, uint(uid))
	if err == userIDRelationSrv.ErrNotFound {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	result, err := userCacheSrv.GetCacheSportResult(db, remoteUID)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}
