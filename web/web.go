package web

import (
	"github.com/gin-gonic/gin"

	apiServer "github.com/inkedawn/JKWXRunner-server/api_server"
	"github.com/inkedawn/JKWXRunner-server/web/controllers"
)

func Run(engine *gin.Engine) error {
	controllers.StartupControllers(engine)
	// old style
	err := apiServer.Run(engine)
	if err != nil {
		return err
	}
	return nil
}
