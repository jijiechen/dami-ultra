package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/yuchanns/kong-exercise-microservices/internal/apis"
	"github.com/yuchanns/kong-exercise-microservices/utils/helpers"
)

func Run(engine *gin.Engine) (err error) {
	err = Register(engine)
	if err != nil {
		return
	}
	return engine.Run(":8080")
}

func Register(engine *gin.Engine) (err error) {
	engine.RedirectTrailingSlash = true

	v1 := engine.Group("/api")
	svc := apis.NewService()
	v1.POST("/message", helpers.BuildHandler(svc.PostMessage))
	return
}
