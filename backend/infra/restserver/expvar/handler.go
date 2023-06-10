package expvar

import (
	"expvar"
	echo "github.com/labstack/echo/v4"
)

func RegisterExpVarHandler(server *echo.Echo) {

	server.GET("/debug/vars", echo.WrapHandler(expvar.Handler()))
}
