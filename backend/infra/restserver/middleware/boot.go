package middleware

import (
	middlewareconfig "code-shooting/infra/restserver/config"
)

type IMiddlewareBootstraper interface {
	Enable(config middlewareconfig.MiddlewareConf) bool
	Order(config middlewareconfig.MiddlewareConf) int
	BindMiddleware(server RestServer, config middlewareconfig.MiddlewareConf)
}
