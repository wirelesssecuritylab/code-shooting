package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"

	"code-shooting/infra/logger"
	middlewareconfig "code-shooting/infra/restserver/config"
)

type NonGetInterceptorBootstraper struct {
	IsMaster     bool
	Availiable   bool
	NonGetConfig *NonGetInterceptorConfig
}

type NonGetInterceptorConfig struct {
	Name      string   `yaml:"name"`
	Order     int      `yaml:"order"`
	Allowlist []string `yaml:"allowlist"`
}

func (boot *NonGetInterceptorBootstraper) Enable(config middlewareconfig.MiddlewareConf) bool {

	_, ok := config["nonGetInterceptor"]
	return ok
}

func (boot *NonGetInterceptorBootstraper) Order(config middlewareconfig.MiddlewareConf) int {

	var order = COMMON_ORDER
	if nonGetInterceptor, ok := config["nonGetInterceptor"]; ok {

		boot.NonGetConfig = &NonGetInterceptorConfig{}

		if err := transformInterfaceToObject(nonGetInterceptor, boot.NonGetConfig); err == nil {
			order = boot.NonGetConfig.Order + order
		}
	}
	return order
}

func (boot *NonGetInterceptorBootstraper) BindMiddleware(server RestServer, config middlewareconfig.MiddlewareConf) {
	boot.NonGetConfig = &NonGetInterceptorConfig{}
	nonGetInterceptor, ok := config["nonGetInterceptor"]
	if !ok {
		logger.Info("the NonGetInterceptor Middleware is not configured, ignore")
		return
	}
	logger.Infof("enable the NonGetInterceptor Middleware, config is %v", nonGetInterceptor)

	err := transformInterfaceToObject(nonGetInterceptor, boot.NonGetConfig)
	if err != nil {
		logger.Infof("the NonGetInterceptor Middleware is error config: %s, ignore", err.Error())
		return
	}

	server.Pre(boot.Handle)

}

func (boot *NonGetInterceptorBootstraper) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// Master状态，不拦截请求
		if boot.IsMaster {
			return next(c)
		}

		//GET请求，不拦截
		method := c.Request().Method
		if method == "GET" {
			return next(c)
		}

		// 白名单中的请求，不拦截
		path := c.Request().URL.Path
		for _, url := range boot.NonGetConfig.Allowlist {
			if ok, _ := regexp.MatchString(url, path); ok {
				return next(c)
			}
		}

		languageOption := strings.ToLower(c.Request().Header.Get("language-option"))
		if strings.Compare("zh-cn", languageOption) == 0 {
			return &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "服务器处于备机状态，不允许该操作。",
			}
		}
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "This operation is not allowed because the server is running in the standby state.",
		}
	}
}
