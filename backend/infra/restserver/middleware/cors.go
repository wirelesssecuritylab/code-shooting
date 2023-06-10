package middleware

import (
	"code-shooting/infra/logger"
	middlewareconfig "code-shooting/infra/restserver/config"

	m "github.com/labstack/echo/v4/middleware"
)

type CORSMiddlewareBootstraper struct {
	CORSConfig *CORSMiddlewareConfig
}

type CORSMiddlewareConfig struct {
	//跨域资源共享
	Name         string           `yaml:"name"`
	m.CORSConfig `yaml:",inline"` // 访问控制，使跨站数据传输更安全
	Order        int              `yaml:"order"`
}

var _ IMiddlewareBootstraper = &CORSMiddlewareBootstraper{}

func (boot *CORSMiddlewareBootstraper) Enable(config middlewareconfig.MiddlewareConf) bool {

	_, ok := config["cors"]
	return ok
}

func (boot *CORSMiddlewareBootstraper) Order(config middlewareconfig.MiddlewareConf) int {

	var order = COMMON_ORDER
	if cors, ok := config["cors"]; ok {

		boot.CORSConfig = &CORSMiddlewareConfig{}

		if err := transformInterfaceToObject(cors, boot.CORSConfig); err == nil {
			order = boot.CORSConfig.Order + order
		}

	}
	return order
}

func (boot *CORSMiddlewareBootstraper) BindMiddleware(server RestServer, config middlewareconfig.MiddlewareConf) {

	CORS, ok := config["cors"]
	if !ok {
		//用户没有配置，按默认处理
		logger.Info("the CORS Middleware is not configured, ignore")
		return
	}
	logger.Infof("enable the CORS Middleware, config is %v", CORS)

	boot.CORSConfig = &CORSMiddlewareConfig{}

	if err := transformInterfaceToObject(CORS, boot.CORSConfig); err != nil {
		//用户配置异常，按默认处理
		logger.Infof("the CORS Middleware is error config: %s, use default CORS Middleware", err.Error())
		server.Use(m.CORS())
		return
	}

	server.Use(m.CORSWithConfig(boot.CORSConfig.CORSConfig))
}
