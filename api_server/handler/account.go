package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv"
	"net/http"
)

func ListAccounts(ctx *gin.Context) {
	db := database.GetDB()
	n, err := accountSrv.CountAccounts(db)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	acc, err := accountSrv.ListAccounts(db, 0, n)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, acc)
}
