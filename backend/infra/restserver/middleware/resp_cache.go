package middleware

import (
	"github.com/labstack/echo/v4"

	"code-shooting/infra/logger"
	middlewareconfig "code-shooting/infra/restserver/config"
)

type RespCacheMiddlewareBootstraper struct {
	RespCacheConfig *RespCacheMiddlewareConfig
}

type RespCacheMiddlewareConfig struct {
	Name                string `yaml:"name"`
	Order               int    `yaml:"order"`
	CacheControl        string `yaml:"cachecontrol"`
	Expires             string `yaml:"expires"`
	Pragma              string `yaml:"pragma"`
	XContentTypeOptions string `yaml:"xcontenttypeoptions"`
}

var _ IMiddlewareBootstraper = &RespCacheMiddlewareBootstraper{}

func (boot *RespCacheMiddlewareBootstraper) Enable(config middlewareconfig.MiddlewareConf) bool {

	_, ok := config["respcache"]
	return ok
}

func (boot *RespCacheMiddlewareBootstraper) Order(config middlewareconfig.MiddlewareConf) int {

	var order = COMMON_ORDER
	if respCache, ok := config["respcache"]; ok {

		boot.RespCacheConfig = &RespCacheMiddlewareConfig{}

		if err := transformInterfaceToObject(respCache, boot.RespCacheConfig); err == nil {
			order = boot.RespCacheConfig.Order + order
		}
	}
	return order
}

func (boot *RespCacheMiddlewareBootstraper) BindMiddleware(server RestServer, config middlewareconfig.MiddlewareConf) {

	boot.RespCacheConfig = &RespCacheMiddlewareConfig{}

	respcache, ok := config["respcache"]
	if !ok {
		logger.Info("the RespCache Middleware is not configured, ignore")
		return
	}
	logger.Infof("enable the RespCache Middleware, config is %v", respcache)

	if err := transformInterfaceToObject(respcache, boot.RespCacheConfig); err != nil {
		logger.Infof("the RespCache Middleware is error config: %s, ignore", err.Error())
		return
	}

	server.Use(boot.RespCacheConfig.getMiddlewareFunc())

}

func (config *RespCacheMiddlewareConfig) getMiddlewareFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			resp := c.Response()

			//Cache-Control：默认值no-cache,no-store
			if len(config.CacheControl) == 0 {
				resp.Header().Set("Cache-Control", "no-cache,no-store")
			} else {
				resp.Header().Set("Cache-Control", config.CacheControl)
			}

			//Expires：默认值0
			if len(config.Expires) == 0 {
				resp.Header().Set("Expires", "0")
			} else {
				resp.Header().Set("Expires", config.Expires)
			}

			//Pragma：默认值no-cache
			if len(config.Pragma) == 0 {
				resp.Header().Set("Pragma", "no-cache")
			} else {
				resp.Header().Set("Pragma", config.Pragma)
			}

			//x-content-type-options：默认值nosniff
			if len(config.XContentTypeOptions) == 0 {
				resp.Header().Set(echo.HeaderXContentTypeOptions, "nosniff")
			} else {
				resp.Header().Set(echo.HeaderXContentTypeOptions, config.XContentTypeOptions)
			}
			resp.Header().Set(echo.HeaderXXSSProtection, "1;mode=block")
			return next(c)
		}
	}
}
