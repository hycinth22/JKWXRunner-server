package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
	"github.com/inkedawn/JKWXRunner-server/web/viewmodels"
)

type AccountRouter struct{}

func (AccountRouter) RegisterToRouter(router gin.IRouter) {
	router.GET("/account", func(context *gin.Context) {
		_, hideTerminated := context.GetQuery("hideTerminated")
		dbSrv := service.NewCommonService()
		leaperSrv := service.NewAccountServiceUpon(dbSrv)
		sportSrv := service.NewUserSportResultServiceUpon(dbSrv)
		var (
			accList []datamodels.Account
			err     error
		)
		if hideTerminated {
			accList, err = leaperSrv.ListAccountsExceptStatus(service.AccountStatusTerminated)
		} else {
			accList, err = leaperSrv.ListAccounts()
		}
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			return
		}
		var resp []*viewmodels.Account
		for _, acc := range accList {
			sport, err := sportSrv.GetLocalUserCacheSportResult(acc.ID)
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
	router.PUT("/account/:id/status", func(ctx *gin.Context) {
		var (
			payload struct {
				Status string
			}
			err error
		)
		tid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		id := uint(tid)
		err = ctx.Bind(&payload)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		accSrv := service.NewAccountService()
		err = accSrv.UpdateAccountStatus(id, payload.Status)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Status(http.StatusAccepted)
	})
	router.POST("/account/:id/finishNow", func(ctx *gin.Context) {
		tid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		id := uint(tid)
		accSrv := service.NewAccountService()
		err = accSrv.FinishAheadOfSchedule(id)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Status(http.StatusOK)
	})
}
