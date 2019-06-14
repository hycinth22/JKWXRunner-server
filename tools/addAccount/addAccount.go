package main

import (
	"fmt"
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/service"
	"github.com/inkedawn/JKWXFucker-server/service/accountSrv"
	"github.com/inkedawn/JKWXFucker-server/service/deviceSrv"
	"github.com/inkedawn/JKWXFucker-server/service/userCacheSrv"
	"github.com/inkedawn/JKWXFucker-server/service/userIDRelationSrv"
	"github.com/inkedawn/go-sunshinemotion"
	"os"
	"strconv"
	"time"
)

var (
	Arg_SchoolID    int64
	Arg_StuNum      string
	Arg_Password    string
	Arg_RunDistance float64
	Arg_Status      string
	Arg_Memo        string
)

func mustParseArgs() {
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
		Arg_RunDistance, err = strconv.ParseFloat(os.Args[4], 64)
		if err != nil {
			panic(err)
		}
	}
	if len(os.Args) >= 6 {
		Arg_Status = os.Args[5]
	}
	if Arg_Status == "" {
		Arg_Status = accountSrv.StatusNormal
	}
	if len(os.Args) >= 7 {
		Arg_Memo = os.Args[6]
	}
}

func main() {
	mustParseArgs()
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

	dev := deviceSrv.FromSSMTDevice(*ssmtDevice)
	err = deviceSrv.SaveDevice(tx, &dev)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Device %d: %+v", dev.ID, dev)
	fmt.Println()
	acc := &accountSrv.Account{
		SchoolID:    Arg_SchoolID,
		StuNum:      Arg_StuNum,
		Password:    Arg_Password,
		RunDistance: Arg_RunDistance,
		DeviceID:    dev.ID,
		Status:      Arg_Status,
		Memo:        Arg_Memo,
	}
	if acc.RunDistance == 0.0 {
		limit := ssmt.GetDefaultLimitParams(info.Sex)
		acc.RunDistance = limit.LimitTotalDistance.Max
	}
	acc.RunDistance = ssmt.NormalizeDistance(acc.RunDistance)
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
}
