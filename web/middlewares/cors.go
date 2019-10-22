package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	context.Header("Access-Control-Allow-Headers", "content-type")
	if context.Request.Method == http.MethodOptions {
		context.AbortWithStatus(http.StatusOK)
	}
	context.Next()
}
