package main

import (
	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/JKWXRunner-server/database"
	"github.com/inkedawn/JKWXRunner-server/service/deviceSrv"
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
