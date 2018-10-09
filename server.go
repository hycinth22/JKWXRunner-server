package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func init() {
	_, err := os.Stat("debug")
	if os.IsNotExist(err) {
		gin.SetMode(gin.ReleaseMode)
	}else{
		gin.SetMode(gin.DebugMode)
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
	engine.Run(":8080")
}
