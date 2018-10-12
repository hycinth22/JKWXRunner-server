package main

import (
	"github.com/gin-gonic/gin"
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
	go runAPIServer()
	go runStaticFileServer()
	select {}
}

func runStaticFileServer() {
	// static files
	engine := gin.Default()
	engine.Use(authMiddleWare)
	engine.Static("/", `./html`)
	engine.Run(":80")
}

func runAPIServer() {
	engine := gin.Default()
	registerCORSRoute(engine)
	registerTicketRoute(engine)
	registerRemoteProfileRoute(engine)
	engine.Run(":8080")
}
