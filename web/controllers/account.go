package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/service"
)

type AccountRouter struct{}

func (AccountRouter) RegisterToRouter(router gin.IRouter) {
	router.GET("/account", func(context *gin.Context) {
		leaperSrv := service.NewAccountService()
		result, err := leaperSrv.ListAccounts()
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			return
		}
		context.JSON(http.StatusOK, result)
	})
	router.POST("/account", notImplementedHandler)
	router.PUT("/account", notImplementedHandler)
	router.DELETE("/account/:id", notImplementedHandler)
}
