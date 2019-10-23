package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service"
	"github.com/inkedawn/JKWXRunner-server/service/accountSrv"
	"github.com/inkedawn/JKWXRunner-server/service/deviceSrv"
	"github.com/inkedawn/JKWXRunner-server/service/userCacheSrv"
	"github.com/inkedawn/JKWXRunner-server/service/userIDRelationSrv"

	"github.com/inkedawn/go-sunshinemotion/v3"
)

var (
	Arg_SchoolID int64
	Arg_StuNum   string
	Arg_Password string
	Arg_OwnerID  int
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
		Arg_OwnerID, err = strconv.Atoi(os.Args[4])
		if err != nil {
			panic(err)
		}
	} else {
		Arg_OwnerID = defaultOwnerID
	}
}

func main() {
	mustParseArgs()
	accSrv := service.NewAccountService()
	acc, err := accSrv.GetAccountByStuNum(Arg_SchoolID, Arg_StuNum)
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
	tx := database.GetDB().Begin()
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
			tx.Rollback()
		} else {
			fmt.Println("Confirm? Enter to continue...")
			_, _ = fmt.Scanln()
			tx.Commit()
		}
	}()

	ssmtDevice := ssmt.GenerateDevice()
	session := new(ssmt.Session)
	session.Device = ssmtDevice
	info, err := session.Login(Arg_SchoolID, Arg_StuNum, "123", ssmt.PasswordHash(Arg_Password))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account Info: %+v", info)
	fmt.Println()
	if info.UserRoleID == userCacheSrv.UserRole_Cheater {
		fmt.Println("!!![WARNING]!!! The User Has been marked as a cheater!")
	}

	fetchTime := time.Now()
	sport, err := session.GetSportResult()
	if err != nil {
		panic(err)
	}
	err = userCacheSrv.SaveCacheSportResult(tx, userCacheSrv.FromSSMTSportResult(*sport, session.User.UserID, fetchTime))
	if err != nil {
		panic(err)
	}
	fmt.Printf("SportResult: %+v", sport)
	fmt.Println()

	dev := deviceSrv.FromSSMTDevice(*ssmtDevice)
	err = deviceSrv.SaveDevice(tx, &dev)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Device %d: %+v", dev.ID, dev)
	fmt.Println()
	limit := ssmt.GetDefaultLimitParams(info.Sex)
	acc = &accountSrv.Account{
		OwnerID:          Arg_OwnerID,
		SchoolID:         Arg_SchoolID,
		StuNum:           Arg_StuNum,
		Password:         Arg_Password,
		RunDistance:      limit.LimitTotalMaxDistance,
		DeviceID:         dev.ID,
		Status:           accountSrv.StatusNormal,
		Memo:             "",
		CheckCheatMarked: true,
	}
	acc.RunDistance = ssmt.NormalizeDistance(acc.RunDistance)
	acc.StartDistance = sport.ActualDistance
	acc.FinishDistance = sport.QualifiedDistance

	err = accountSrv.SaveAccount(tx, acc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account %d: %+v", acc.ID, acc)
	fmt.Println()
	err = userIDRelationSrv.SaveRelation(tx, acc.ID, session.User.UserID)
	if err != nil {
		panic(err)
	}

	if info.UserRoleID == userCacheSrv.UserRole_Cheater {
		fmt.Println("!!![WARNING]!!! Disable CheckCheatMarked! Confirm?")
		fmt.Println("Enter to continue...")

		_, _ = fmt.Scanln()
	}
}
