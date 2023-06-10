package middleware

import (
	m "github.com/labstack/echo/v4/middleware"

	"code-shooting/infra/logger"
	middlewareconfig "code-shooting/infra/restserver/config"
)

type RequestIDMiddlewareBootstraper struct {
	RequestIDConfig *RequestIDMiddlewareConfig
}

type RequestIDMiddlewareConfig struct {
	// 为请求生成唯一的ID
	Name              string `yaml:"name"`
	m.RequestIDConfig `yaml:",inline"`
	Order             int `yaml:"order"`
}

var _ IMiddlewareBootstraper = &RequestIDMiddlewareBootstraper{}

func (boot *RequestIDMiddlewareBootstraper) Enable(config middlewareconfig.MiddlewareConf) bool {

	_, ok := config["requestid"]
	return ok
}

func (boot *RequestIDMiddlewareBootstraper) Order(config middlewareconfig.MiddlewareConf) int {
	var order = COMMON_ORDER
	if requestID, ok := config["requestid"]; ok {

		boot.RequestIDConfig = &RequestIDMiddlewareConfig{}

		if err := transformInterfaceToObject(requestID, boot.RequestIDConfig); err == nil {
			order = boot.RequestIDConfig.Order + order
		}
	}
	return order
}

func (boot *RequestIDMiddlewareBootstraper) BindMiddleware(server RestServer, config middlewareconfig.MiddlewareConf) {

	boot.RequestIDConfig = &RequestIDMiddlewareConfig{}

	requestID, ok := config["requestid"]
	if !ok {
		//用户没有配置
		logger.Infof("the RequestID Middleware is not configured, ignore")
		return
	}
	logger.Infof("enable the RequestID Middleware, config is %v", requestID)

	if err := transformInterfaceToObject(requestID, boot.RequestIDConfig); err != nil {
		logger.Info("the RequestID Middleware is error config, use default RequestID Middleware")
		server.Use(m.RequestID())
		return
	}

	server.Use(m.RequestIDWithConfig(boot.RequestIDConfig.RequestIDConfig))

}
