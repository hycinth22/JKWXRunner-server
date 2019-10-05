package apiServer

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/api_server/handler"
	"github.com/inkedawn/JKWXRunner-server/service"
)

var cors = false

func init() {
	RequestRoutersTable = []RequestRouter{
		// CORSRouter, // default disabled by security cause. Only enable it if you need CORS actually.
		AccountRouter,
		UserInfoRouter,
		SportResultRouter,
		ErrorsCollect,
	}
	if cors {
		RequestRoutersTable = append(RequestRoutersTable, CORSRouter)
	}
}

func CORSRouter(router gin.IRouter) {
	router.Use(func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")

	})
	router.OPTIONS("/:all", func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
		context.Header("Access-Control-Allow-Headers", "*")
	})
}

func ErrorsCollect(router gin.IRouter) {
	fmt.Println("s0.0")
	router.Use(func(context *gin.Context) {
		fmt.Println("s0")
		context.Next()
		fmt.Println("s1")
		for _, err := range context.Errors {
			fmt.Println("s2")
			if service.IsInternalError(err) {
				log.Println(*context.Request, err, service.UnwrapInternalError(err))
			} else {
				log.Println(*context.Request, err)
			}
		}
	})
}

func AccountRouter(router gin.IRouter) {
	router.GET("/account", handler.ListAccounts)
	router.POST("/account", handler.NotImplemented)
	router.PUT("/account", handler.NotImplemented)
	router.DELETE("/account/:id", handler.NotImplemented)
	router.GET("/account/:id/logs", handler.ListLogsByID)
	router.PUT("/account/:id/status", handler.UpdateAccountStatus)
}

func UserInfoRouter(router gin.IRouter) {
	router.GET("/userInfo/:username", handler.NotImplemented)
}

func SportResultRouter(router gin.IRouter) {
	router.GET("/sportResult/:uid", handler.QueryCacheUserSportResult)
}
