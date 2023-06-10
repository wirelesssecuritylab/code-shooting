package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"

	middlewareconfig "code-shooting/infra/restserver/config"
)

const (
	X_HTTP_METHOD          = "X-HTTP-Method"
	X_HTTP_METHOD_OVERRIDE = "X-HTTP-Method-Override"
	X_METHOD_OVERRIDE      = "X-Method-Override"
)

type MethodOverrideMiddlewareBootstraper struct {
	MethodOverrideConfig *MethodOverrideMiddlewareConfig
}

type MethodOverrideMiddlewareConfig struct {
	Name  string `yaml:"name"`
	Order int    `yaml:"order"`
}

var _ IMiddlewareBootstraper = &MethodOverrideMiddlewareBootstraper{}

func (boot *MethodOverrideMiddlewareBootstraper) Enable(config middlewareconfig.MiddlewareConf) bool {

	_, ok := config["methodOverride"]
	return ok
}

func (boot *MethodOverrideMiddlewareBootstraper) Order(config middlewareconfig.MiddlewareConf) int {

	var order = COMMON_ORDER
	if header, ok := config["methodOverride"]; ok {

		boot.MethodOverrideConfig = &MethodOverrideMiddlewareConfig{}

		if err := transformInterfaceToObject(header, boot.MethodOverrideConfig); err == nil {
			order = boot.MethodOverrideConfig.Order + order
		}
	}
	return order
}

func (boot *MethodOverrideMiddlewareBootstraper) BindMiddleware(server RestServer, config middlewareconfig.MiddlewareConf) {
	server.Pre(MethodOverrideWithConfig())
}

// RequestIDWithConfig returns a X_HTTP_METHOD middleware with config.
func MethodOverrideWithConfig() echo.MiddlewareFunc {
	// Defaults

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			//获取Header
			header := c.Request().Header
			url := c.Request().URL
			if (header != nil && (header.Get(X_HTTP_METHOD) != "" ||
				header.Get(X_HTTP_METHOD_OVERRIDE) != "" ||
				header.Get(X_METHOD_OVERRIDE) != "")) ||
				strings.Contains(url.Path, "_method") ||
				strings.Contains(url.RawQuery, "_method") {
				return echo.ErrMethodNotAllowed
			}
			return next(c)
		}
	}
}
