package middleware

import "github.com/labstack/echo/v4"

type IMiddleware interface {
	GetName() string
	SetConfig(params interface{}) error
	Handle(next echo.HandlerFunc) echo.HandlerFunc
}
