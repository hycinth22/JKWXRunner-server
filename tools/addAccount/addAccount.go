package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/datamodels"
	"github.com/inkedawn/JKWXRunner-server/service"
)

var (
	Arg_SchoolID       int64
	Arg_StuNum         string
	Arg_Password       string
	Arg_FinishDistance float64
	Arg_OwnerID        int
)

func mustParseArgs() {
	const defaultOwnerID = 1
	var err error
	if len(os.Args) < 4 {
		panic("too few arguments")
	}
	Arg_SchoolID, err = strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		panic(err)
	}
	Arg_StuNum = os.Args[2]
	Arg_Password = os.Args[3]
	if len(os.Args) >= 5 {
		Arg_FinishDistance, _ = strconv.ParseFloat(os.Args[4], 64)
	}
	if len(os.Args) >= 6 {
		Arg_OwnerID, err = strconv.Atoi(os.Args[5])
		if err != nil {
			panic(err)
		}
	} else {
		Arg_OwnerID = defaultOwnerID
	}
}

func main() {
	mustParseArgs()
	common := service.NewCommonService()
	common.Begin()
	defer func() {
		x := recover()

		// unwrap error if it's a internal error
		if xx, ok := x.(error); ok {
			if service.IsInternalError(xx) {
				x = service.UnwrapInternalError(xx)
			}
		}

		if x != nil {
			fmt.Println(x)
			common.Rollback()
		} else {
			fmt.Println("Confirm? Enter to continue...")
			_, _ = fmt.Scanln()
			common.Commit()
		}
	}()
	accSrv := service.NewAccountServiceUpon(common)
	acc, err := accSrv.GetActiveAccountByStuNum(Arg_SchoolID, Arg_StuNum)
	switch err {
	case service.ErrNoAccount:
		break
	case nil:
		if acc != nil {
			fmt.Println("帐号已存在. 状态是：", acc.Status)
			return
		}
	default:
		panic(err)
	}
	ssmtDevice := ssmt.GenerateDevice()
	session := new(ssmt.Session)
	session.Device = ssmtDevice
	info, err := session.Login(Arg_SchoolID, Arg_StuNum, "123", ssmt.PasswordHash(Arg_Password))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account Info: %+v", info)
	fmt.Println()
	if info.UserRoleID == service.UserRole_Cheater {
		fmt.Println("!!![WARNING]!!! The User Has been marked as a cheater!")
	}

	fetchTime := time.Now()
	sport, err := session.GetSportResult()
	if err != nil {
		panic(err)
	}

	err = service.NewUserSportResultServiceUpon(common).SaveCacheSportResult(datamodels.CacheUserSportResultFromSSMTSportResult(*sport, session.User.UserID, fetchTime))
	if err != nil {
		panic(err)
	}
	fmt.Printf("SportResult: %+v", sport)
	fmt.Println()

	dev := datamodels.DeviceFromSSMTDevice(*ssmtDevice)
	err = service.NewDeviceServiceUpon(common).SaveDevice(&dev)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Device %d: %+v", dev.ID, dev)
	fmt.Println()
	limit := ssmt.GetDefaultLimitParams(info.SchoolID, info.Sex)
	acc = &datamodels.Account{
		OwnerID:          Arg_OwnerID,
		SchoolID:         Arg_SchoolID,
		StuNum:           Arg_StuNum,
		Password:         Arg_Password,
		RunDistance:      limit.LimitTotalMaxDistance,
		DeviceID:         dev.ID,
		Status:           service.AccountStatusNormal,
		Memo:             "",
		CheckCheatMarked: sql.NullBool{Valid: false},
	}
	acc.RunDistance = ssmt.NormalizeDistance(acc.RunDistance)
	acc.StartDistance = sport.ActualDistance
	acc.FinishDistance = sport.QualifiedDistance
	if Arg_FinishDistance != 0 {
		acc.FinishDistance = Arg_FinishDistance
	}

	err = accSrv.SaveAccount(acc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account %d: %+v", acc.ID, acc)
	fmt.Println()
	err = service.NewUserIDRelServiceUpon(common).SaveRelation(acc.ID, session.User.UserID)
	if err != nil {
		panic(err)
	}

	if info.UserRoleID == service.UserRole_Cheater {
		fmt.Println("!!![WARNING]!!! Disable CheckCheatMarked! Confirm?")
		_, _ = fmt.Scanln()
		err = accSrv.SetCheckCheaterFlag(acc.ID, false)
		if err != nil {
			fmt.Println("[ERROR] an error occurred when set flag. ", err)
		}
	}
}
