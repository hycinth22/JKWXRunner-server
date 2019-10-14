package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service"
	"github.com/inkedawn/JKWXRunner-server/service/userCacheSrv"
	"github.com/inkedawn/JKWXRunner-server/viewFormat"
)

type AccountRouter struct{}

type account struct {
	ID               uint
	CreatedAt        string
	SchoolID         int64
	StuNum           string
	Memo             string
	Status           string
	RunDistance      float64
	StartDistance    float64
	FinishDistance   float64
	CurrentDistance  float64
	CheckCheatMarked bool
	LastResult       string
	LastTime         string
}

func (AccountRouter) RegisterToRouter(router gin.IRouter) {
	router.GET("/account", func(context *gin.Context) {
		leaperSrv := service.NewAccountService()
		accList, err := leaperSrv.ListAccounts()
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			return
		}
		var resp []account
		for _, acc := range accList {
			current := -0.0
			sport, err := userCacheSrv.GetLocalUserCacheSportResult(database.GetDB(), acc.ID)
			if err == nil {
				current = sport.ComputedDistance
			}
			resp = append(resp, account{
				ID:               acc.ID,
				CreatedAt:        viewFormat.TimeFormat(acc.CreatedAt),
				SchoolID:         acc.SchoolID,
				StuNum:           acc.StuNum,
				Memo:             acc.Memo,
				Status:           acc.Status,
				RunDistance:      acc.RunDistance,
				StartDistance:    acc.StartDistance,
				FinishDistance:   acc.FinishDistance,
				CurrentDistance:  current,
				CheckCheatMarked: acc.CheckCheatMarked,
				LastResult:       acc.LastResult,
				LastTime:         viewFormat.TimeFormat(acc.LastTime),
			})
		}
		context.JSON(http.StatusOK, resp)
	})
	router.POST("/account", notImplementedHandler)
	router.PUT("/account", notImplementedHandler)
	router.DELETE("/account/:id", notImplementedHandler)
}
