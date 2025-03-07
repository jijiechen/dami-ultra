package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/jijiechen/dami-ultra/internal/apis"
	"github.com/jijiechen/dami-ultra/utils/helpers"
	"os"
)

var V2APIEnabled = os.Getenv("V2_API_ENABLED") == "true"
var ServerListenAddr = os.Getenv("SERVER_ADDRESS")

func Run(engine *gin.Engine) (err error) {
	err = Register(engine)
	if err != nil {
		return
	}

	listenAddr := ServerListenAddr
	if listenAddr == "" {
		listenAddr = ":8080"
	}
	return engine.Run(listenAddr)
}

func Register(engine *gin.Engine) (err error) {
	engine.RedirectTrailingSlash = true

	apiEndpoints := engine.Group("/api")
	svc := apis.NewService()
	if V2APIEnabled {
		apiEndpoints.POST("/message", helpers.BuildHandler(svc.PostOperationMessage))
	} else {
		apiEndpoints.POST("/message", helpers.BuildHandler(svc.PostMessages))
	}

	apiEndpoints.POST("/message/v1", helpers.BuildHandler(svc.PostMessages))
	apiEndpoints.POST("/message/v2", helpers.BuildHandler(svc.PostOperationMessage))
	return
}
