package middlewares

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/inkedawn/JKWXRunner-server/service"
)

func ErrorCollect(context *gin.Context) {
	context.Next()
	for _, err := range context.Errors {
		var interErr *service.InternalError
		if errors.As(err, &interErr) {
			log.Println("ErrorCollect: ", *context.Request, "\n", interErr.Unwrap().Error())
		} else {
			log.Println("ErrorCollect: ", *context.Request, "\n", err.Error())
		}
	}
}
