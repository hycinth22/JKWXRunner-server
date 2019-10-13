package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/config"
	"github.com/inkedawn/JKWXRunner-server/web"
)

// inject when build
// go build/run -ldflags "-X 'main.lastCommit=hello lastCommit' -X 'main.lastCommitTime=hello lastCommitTime'"
var (
	lastCommit     = "not set(maybe test version)"
	lastCommitTime = "not set(maybe test version)"
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
	fmt.Println("BuildVersion:", lastCommit)
	fmt.Println("BuildVersionTime:", lastCommitTime)
	err := web.Run(gin.New())
	if err != nil {
		log.Fatal(err)
	}
}
