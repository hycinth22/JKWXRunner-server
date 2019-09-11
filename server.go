package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/api_server"
	"github.com/inkedawn/JKWXRunner-server/config"
)

func init() {
	//noinspection GoBoolExpressions
	if config.Release {
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
	err := apiServer.Run(gin.New())
	if err != nil {
		log.Fatal(err)
	}
}
