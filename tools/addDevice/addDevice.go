package main

import (
	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service/deviceSrv"
	"github.com/inkedawn/go-sunshinemotion"
)

func main() {
	db := database.GetDB()

	dev := deviceSrv.FromSSMTDevice(*ssmt.GenerateDevice())
	err := deviceSrv.SaveDevice(db, &dev)
	if err != nil {
		println(err.Error())
	}
	println(dev.ID)
}
