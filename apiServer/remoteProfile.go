package apiServer

import (
	"github.com/gin-gonic/gin"
	sunshinemotion "github.com/inkedawn/go-sunshinemotion"
	"log"
	"net/http"
	"strconv"
	"github.com/inkedawn/JKWXFucker-server/model"
)

type RemoteProfile struct {
	StudentName   string `json:"studentName"`
	StudentNumber string `json:"studentNumber"`
	Sex           string `json:"sex"`
	SportResult
}

type SportResult struct {
	CompletedDistance float64 `json:"completedDistance"`
	QualifiedDistance float64 `json:"qualifiedDistance"`
}

func registerRemoteProfileRoute(router gin.IRouter) {
	router.GET("/remoteProfile/:schoolID/:username", getRemoteProfile)
}

func getRemoteProfile(context *gin.Context) {
	schoolID, err := strconv.ParseInt(context.Param("schoolID"), 10, 64)
	if err != nil {
		context.Error(err)
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	username := context.Param("username")
	password := context.Query("password")
	s, err := model.GetSessionByUsername(username)
	if err != nil {
		if err == model.ErrSessionNotFound {
			if password == "" {
				context.String(200, "需要登录信息。")
				return
			}
			s = sunshinemotion.CreateSession()
			err = s.LoginEx(username, "123", sunshinemotion.PasswordHash(password), schoolID)
			if err != nil {
				context.Error(err)
				context.String(200, "登录错误：%s", err)
				return
			}
		}else {
			context.Error(err)
			context.String(503, "获取SessionStore错误：%s", err)
			return
		}
	}
	log.Println("Get RemoteProfile, user", username, s)
	profile := RemoteProfile{
		StudentName:   s.UserInfo.StudentName,
		StudentNumber: s.UserInfo.StudentNumber,
		Sex:           s.UserInfo.Sex,
	}
	r, err := s.GetSportResult()
	if err == sunshinemotion.ErrTokenExpired {
		err = s.LoginEx(username, "123", sunshinemotion.PasswordHash(password), schoolID)
		if err != nil {
			context.Error(err)
			context.String(200, "登录错误：%s", err)
			return
		}
		// retry
		getRemoteProfile(context)
		return
		
	}
	if err != nil {
		context.Error(err)
	} else {
		profile.SportResult = SportResult{
			CompletedDistance: r.Distance,
			QualifiedDistance: r.Qualified,
		}
	}
	context.JSON(http.StatusOK, profile)
}