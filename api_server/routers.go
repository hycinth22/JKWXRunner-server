package apiServer

import (
	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/api_server/handler"
)

func init() {
	RequestRoutersTable = []RequestRouter{
		AccountRouter,
		UserInfoRouter,
		SportResultRouter,
	}
}

func AccountRouter(router gin.IRouter) {
	//router.GET("/account", handler.ListAccounts)
	//router.POST("/account", handler.NotImplemented)
	//router.PUT("/account", handler.NotImplemented)
	//router.DELETE("/account/:id", handler.NotImplemented)
	router.GET("/account/:id/logs", handler.ListLogsByID)
	//router.PUT("/account/:id/status", handler.UpdateAccountStatus)
}

func UserInfoRouter(router gin.IRouter) {
	router.GET("/userInfo/:username", handler.NotImplemented)
}

func SportResultRouter(router gin.IRouter) {
	router.GET("/sportResult/:uid", handler.QueryCacheUserSportResult)
}
