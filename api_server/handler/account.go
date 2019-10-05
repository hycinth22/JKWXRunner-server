package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv"
	"github.com/inkedawn/JKWXRunner-server/service/userCacheSrv"
	"github.com/inkedawn/JKWXRunner-server/viewFormat"
)

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
	ctx.JSON(http.StatusOK, resp)
}

func UpdateAccountStatus(ctx *gin.Context) {
	var (
		payload struct {
			Status string
		}
		err error
	)
	db := database.GetDB()
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
	acc, err := accountSrv.GetAccount(db, id)
	if err == accountSrv.ErrNoAccount {
		ctx.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	acc.Status = payload.Status
	err = accountSrv.SaveAccount(db, acc)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusAccepted)
}
