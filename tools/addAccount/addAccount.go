package main

import (
	"fmt"
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/service"
	"github.com/inkedawn/JKWXFucker-server/service/accountSrv"
	"github.com/inkedawn/JKWXFucker-server/service/deviceSrv"
	"github.com/inkedawn/go-sunshinemotion"
	"os"
	"strconv"
)

var (
	Arg_SchoolID    int64
	Arg_StuNum      string
	Arg_Password    string
	Arg_RunDistance float64
)

func mustParseArgs() {
	var err error
	if len(os.Args) < 5 {
		panic("too few arguments")
	}
	Arg_SchoolID, err = strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		panic(err)
	}
	Arg_StuNum = os.Args[2]
	Arg_Password = os.Args[3]
	Arg_RunDistance, err = strconv.ParseFloat(os.Args[4], 64)
	if err != nil {
		panic(err)
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
	dev := deviceSrv.FromSSMTDevice(*ssmtDevice)
	err := deviceSrv.SaveDevice(tx, &dev)
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
		Status:      accountSrv.StatusSuspend,
	}
	err = accountSrv.SaveAccount(tx, acc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account %d: %+v", acc.ID, acc)
	fmt.Println()
}
