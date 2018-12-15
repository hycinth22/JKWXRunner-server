package main

import (
	"github.com/gin-gonic/gin"
	"github.com/inkedawn/JKWXFucker-server/apiServer"
	"github.com/inkedawn/JKWXFucker-server/staticFileServer"
	"log"
	"os"
)

func init() {
	_, err := os.Stat("debug")
	if os.IsNotExist(err) {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	workDir, err := os.Getwd()
	if err == nil {
		log.Println("Working Directory: ", workDir)
	} else {
		log.Println("Get Working Directory Fail.")
	}
}

func main() {
	engine := gin.Default()
	go apiServer.Run(engine)
	go staticFileServer.Run(gin.Default())
	select {}
}
