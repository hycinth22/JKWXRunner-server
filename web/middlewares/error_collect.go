package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/service"
)

func ErrorCollect(context *gin.Context) {
	context.Next()
	for _, err := range context.Errors {
		if service.IsInternalError(err) {
			log.Println("ErrorCollect: ", *context.Request, err, service.UnwrapInternalError(err))
		} else {
			log.Println("ErrorCollect: ", *context.Request, err)
		}
	}
}
