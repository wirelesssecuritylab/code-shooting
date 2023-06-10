package middleware

import (
	m "github.com/labstack/echo/v4/middleware"

	"code-shooting/infra/logger"
	middlewareconfig "code-shooting/infra/restserver/config"
)

type RecoverMiddlewareBootstraper struct {
	RecoverConfig *RecoverMiddlewareConfig
}

type RecoverMiddlewareConfig struct {
	Name            string `yaml:"name"`
	m.RecoverConfig `yaml:",inline"`
	Order           int `yaml:"order"`
}

var _ IMiddlewareBootstraper = &RecoverMiddlewareBootstraper{}

func (boot *RecoverMiddlewareBootstraper) Enable(config middlewareconfig.MiddlewareConf) bool {

	_, ok := config["recover"]
	return ok
}

func (boot *RecoverMiddlewareBootstraper) Order(config middlewareconfig.MiddlewareConf) int {

	var order = COMMON_ORDER
	if recover, ok := config["recover"]; ok {
		boot.RecoverConfig = &RecoverMiddlewareConfig{}
		if err := transformInterfaceToObject(recover, boot.RecoverConfig); err == nil {
			order = boot.RecoverConfig.Order + order
			return order
		}
	}
	return order - 1
}

func (boot *RecoverMiddlewareBootstraper) BindMiddleware(server RestServer, config middlewareconfig.MiddlewareConf) {

	boot.RecoverConfig = &RecoverMiddlewareConfig{}

	recover, ok := config["recover"]
	if !ok {
		//按默认处理
		logger.Info("the Recover Middleware is not configured, enable default Recover Middleware")
		server.Use(m.Recover())
		return
	}
	logger.Infof("enable the Recover Middleware, config is %v", recover)

	if err := transformInterfaceToObject(recover, boot.RecoverConfig); err != nil {
		logger.Info("the Recover Middleware is error config: %s, use default Recover Middleware", err.Error())
		server.Use(m.Recover())
		return
	}

	server.Use(m.RecoverWithConfig(boot.RecoverConfig.RecoverConfig))
}
