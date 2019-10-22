package middlewares

import (
	"log"
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
)

type middleWare = gin.HandlerFunc

var enabled = []middleWare{CORS, ErrorCollect}

func InjectMiddleWares(router gin.IRouter) {
	for _, m := range enabled {
		log.Println("Enable MiddleWare", runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name())
		router.Use(m)
	}
}
