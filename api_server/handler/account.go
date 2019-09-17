package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv"
	"github.com/inkedawn/JKWXRunner-server/service/userCacheSrv"
)

type account struct {
	ID               uint
	CreatedAt        time.Time
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
	LastTime         time.Time
}

func ListAccounts(ctx *gin.Context) {
	db := database.GetDB()
	n, err := accountSrv.CountAccounts(db)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	accList, err := accountSrv.ListAccounts(db, 0, n)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var (
		resp []account
	)
	for _, acc := range accList {
		current := -0.0
		sport, err := userCacheSrv.GetLocalUserCacheSportResult(db, acc.ID)
		if err == nil {
			current = sport.ComputedDistance
		}
		resp = append(resp, account{
			ID:               acc.ID,
			CreatedAt:        acc.CreatedAt,
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
			LastTime:         acc.LastTime,
		})
	}
	ctx.JSON(http.StatusOK, resp)
}
