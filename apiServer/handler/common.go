package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NotImplemented(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}
