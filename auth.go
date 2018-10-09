package main

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
}

func getTodayAuth() string {
	user := "admin"
	password := time.Now().Format("20060102")
	return authorizationHeader(user, password)
}

// modified https://github.com/gin-gonic/gin/blob/master/auth.go#L67
func authMiddleWare(c *gin.Context) {
	pass := c.GetHeader("Authorization") == getTodayAuth()
	if !pass {
		// Credentials doesn't match, we return 401 and abort handlers chain.
		c.Header("WWW-Authenticate", "Basic realm="+strconv.Quote("adminPages"))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
