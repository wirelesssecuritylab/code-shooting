package middleware

import (
	m "github.com/labstack/echo/v4/middleware"

	"code-shooting/infra/logger"
	middlewareconfig "code-shooting/infra/restserver/config"
)

type BodyLimitMiddlewareBootstraper struct {
	BodyLimitConfig *BodyLimitMiddlewareConfig
}

type BodyLimitMiddlewareConfig struct {
	Name              string `yaml:"name"`
	m.BodyLimitConfig `yaml:",inline"`
	Order             int `yaml:"order"`
}

var _ IMiddlewareBootstraper = &BodyLimitMiddlewareBootstraper{}

func (boot *BodyLimitMiddlewareBootstraper) Enable(config middlewareconfig.MiddlewareConf) bool {

	_, ok := config["bodylimit"]
	return ok
}

func (boot *BodyLimitMiddlewareBootstraper) Order(config middlewareconfig.MiddlewareConf) int {

	var order = COMMON_ORDER
	if bodylimit, ok := config["bodylimit"]; ok {

		boot.BodyLimitConfig = &BodyLimitMiddlewareConfig{}

		if err := transformInterfaceToObject(bodylimit, boot.BodyLimitConfig); err == nil {
			order = boot.BodyLimitConfig.Order + order
		}

	}
	return order
}

func (boot *BodyLimitMiddlewareBootstraper) BindMiddleware(server RestServer, config middlewareconfig.MiddlewareConf) {

	boot.BodyLimitConfig = &BodyLimitMiddlewareConfig{}

	bodylimit, ok := config["bodylimit"]
	if !ok {
		logger.Info("the BodyLimit Middleware is not configured, ignore")
		return
	}
	logger.Infof("enable the BodyLimit Middleware, config is %v", bodylimit)

	err := transformInterfaceToObject(bodylimit, boot.BodyLimitConfig)
	if err != nil {
		logger.Infof("the BodyLimit Middleware is error config: %s, ignore", err.Error())
		return
	}

	server.Use(m.BodyLimitWithConfig(boot.BodyLimitConfig.BodyLimitConfig))

}
