package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/service/accountSrv"
	"net/http"
)

func ListAccounts(ctx *gin.Context) {
	acc, err := accountSrv.ListAccounts(database.GetDB(), 0, 100)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, acc)
}
func AddAccount(ctx *gin.Context) {

}
func UpdateAccount(ctx *gin.Context) {

}
