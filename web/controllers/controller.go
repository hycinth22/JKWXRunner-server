// 路由分发API请求，以及将结果转化为响应。
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var controllers = []Controller{
	AccountRouter{},
}

type Controller interface {
	RegisterToRouter(router gin.IRouter)
}

func StartupControllers(router gin.IRouter) {
	for _, c := range controllers {
		c.RegisterToRouter(router)
	}
}

func notImplementedHandler(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}
