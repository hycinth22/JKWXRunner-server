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
			var current, qualified = -0.0, -0.0
			if err == nil {
				current, qualified = sport.ComputedDistance, sport.QualifiedDistance
			}
			resp = append(resp, viewmodels.AccountView(&acc, current, qualified))
		}
		context.JSON(http.StatusOK, resp)
	})
	router.POST("/account", func(context *gin.Context) {
		var param struct {
			SchoolID int64
			StuNum   string
			Password string
		}
		if err := context.Bind(&param); err != nil {
			context.Error(err)
			return
		}
		srv := service.NewAccountService()
		acc, err := srv.CreateAccount(param.SchoolID, param.StuNum, param.Password)
		if err != nil {
			context.Error(err)
			return
		}
		context.JSON(http.StatusCreated, acc)
	})
	router.PUT("/account", notImplementedHandler)
	router.DELETE("/account/:id", notImplementedHandler)
}
