package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/yuchanns/kong-exercise-microservices/internal/apis"
	"github.com/yuchanns/kong-exercise-microservices/internal/business"
	"github.com/yuchanns/kong-exercise-microservices/internal/repos"
	"github.com/yuchanns/kong-exercise-microservices/utils/helpers"
	"github.com/yuchanns/kong-exercise-microservices/utils/middlewares"
	"go.uber.org/dig"
)

func Run(engine *gin.Engine) (err error) {
	err = Register(engine)
	if err != nil {
		return
	}
	return engine.Run(":8080")
}

func Register(engine *gin.Engine) (err error) {
	engine.Use(middlewares.UseTenantDB())
	engine.RedirectTrailingSlash = true

	c := dig.New()

	err = c.Provide(func() repos.ServiceRepo {
		return repos.NewImplServiceRepo()
	})
	if err != nil {
		return
	}
	err = c.Provide(func(repo repos.ServiceRepo) business.IService {
		return business.NewImplService(repo)
	})
	if err != nil {
		return
	}
	err = c.Provide(func(biz business.IService) *apis.Service {
		return apis.NewService(biz)
	})
	if err != nil {
		return
	}

	v1 := engine.Group("/api/v1")

	err = c.Invoke(func(svc *apis.Service) {
		v1.Group("/service").
			GET("/list", helpers.BuildHandler(svc.List)).
			GET("/:id", helpers.BuildHandlerUri(svc.Get)).
			POST("/", helpers.BuildHandler(svc.Create)).
			PUT("/", helpers.BuildHandler(svc.Update))
	})
	return
}
