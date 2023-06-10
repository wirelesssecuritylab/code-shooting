package middleware

import (
	"github.com/labstack/echo/v4"
)

type RestServer interface {
	BindMiddleware(m IMiddleware)
	Use(middleware ...echo.MiddlewareFunc)
	Pre(middleware ...echo.MiddlewareFunc)
}
