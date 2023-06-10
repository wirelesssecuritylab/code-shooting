package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/tylerb/gls"

	middlewareconfig "code-shooting/infra/restserver/config"
)

const (
	REQUEST_ID_NAME     = "x-request-id"
	TRACE_ID_NAME       = "x-b3-traceid"
	SPAN_ID_NAME        = "x-b3-spanid"
	PARENT_SPAN_ID_NAME = "x-b3-parentspanid"
	SAMPLED_NAME        = "x-b3-sampled"
	FLAGS_NAME          = "x-b3-flags"
	SPAN_CONTEXT_NAME   = "x-ot-span-context"

	X_DEXMESH_PREFIX = "x-dexmesh-"
)

type HeaderMiddlewareBootstraper struct {
	HeaderConfig *HeaderMiddlewareConfig
}

type HeaderMiddlewareConfig struct {
	Name  string `yaml:"name"`
	Order int    `yaml:"order"`
}

var _ IMiddlewareBootstraper = &HeaderMiddlewareBootstraper{}

func (boot *HeaderMiddlewareBootstraper) Enable(config middlewareconfig.MiddlewareConf) bool {

	_, ok := config["header"]
	return ok
}

func (boot *HeaderMiddlewareBootstraper) Order(config middlewareconfig.MiddlewareConf) int {

	var order = COMMON_ORDER
	if header, ok := config["header"]; ok {

		boot.HeaderConfig = &HeaderMiddlewareConfig{}

		if err := transformInterfaceToObject(header, boot.HeaderConfig); err == nil {
			order = boot.HeaderConfig.Order + order
		}
	}
	return order
}

func (boot *HeaderMiddlewareBootstraper) BindMiddleware(server RestServer, config middlewareconfig.MiddlewareConf) {
	server.Use(HeaderWithConfig())
}

// RequestIDWithConfig returns a X-Request-ID middleware with config.
func HeaderWithConfig() echo.MiddlewareFunc {
	// Defaults

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			//获取Header
			header := c.Request().Header
			if header != nil && header.Get(TRACE_ID_NAME) != "" {
				gls.Set("header", header)
				defer gls.Cleanup()
			}
			return next(c)

		}
	}
}
