package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service"
	"github.com/inkedawn/JKWXRunner-server/service/userCacheSrv"
	"github.com/inkedawn/JKWXRunner-server/web/viewmodels"
)

type AccountRouter struct{}

func (AccountRouter) RegisterToRouter(router gin.IRouter) {
	router.GET("/account", func(context *gin.Context) {
		leaperSrv := service.NewAccountService()
		accList, err := leaperSrv.ListAccounts()
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			return
		}
		var resp []*viewmodels.Account
		for _, acc := range accList {
			sport, err := userCacheSrv.GetLocalUserCacheSportResult(database.GetDB(), acc.ID)
			current := 0.0
			if err == nil {
				current = sport.ComputedDistance
			}
			resp = append(resp, viewmodels.NewAccount(&acc, current))
		}
		context.JSON(http.StatusOK, resp)
	})
	router.POST("/account", notImplementedHandler)
	router.PUT("/account", notImplementedHandler)
	router.DELETE("/account/:id", notImplementedHandler)
}
