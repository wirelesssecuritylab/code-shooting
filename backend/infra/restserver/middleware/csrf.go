package middleware

import (
	"bufio"
	"bytes"
	"crypto/subtle"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	c "code-shooting/infra/config"
	"code-shooting/infra/logger"
	middlewareconfig "code-shooting/infra/restserver/config"
)

const (
	Z_EXTENT                  = "Z-EXTENT"
	OAUTH_PASSWORD_CREDENTIAL = "OAUTH2-PASSWORD-CREDENTIAL"
	OAUTH_CLIENT_LOGIN        = "/api/oauth2/v1/login"
	OAUTH_CLIENT_LOGOUT       = "/api/oauth2/v1/logout"
	USER_ADMIN                = "admin"
	DEFAULT_FILTER_SUFFIX     = ".html,.css,.js,.png,.properties,.gif,.jpg,.ttf,.woff,.ico,.svg,.json"
	CSRFTOKEN                 = "forgerydefense"
)

type CSRFMiddlewareBootstraper struct {
	CSRFConfig *CSRFMiddlewareConfig
}

type CSRFMiddlewareConfig struct {
	// 跨站请求伪造
	// Enable bool `yaml:"enable"`
	Name  string `yaml:"name"`
	Order int    `yaml:"order"`
}

var _ IMiddlewareBootstraper = &CSRFMiddlewareBootstraper{}

func (boot *CSRFMiddlewareBootstraper) Enable(config middlewareconfig.MiddlewareConf) bool {

	_, ok := config["csrf"]
	return ok
}

func (boot *CSRFMiddlewareBootstraper) Order(config middlewareconfig.MiddlewareConf) int {

	var order = COMMON_ORDER
	if csrf, ok := config["csrf"]; ok {

		boot.CSRFConfig = &CSRFMiddlewareConfig{}

		if err := transformInterfaceToObject(csrf, boot.CSRFConfig); err == nil {
			order = boot.CSRFConfig.Order + order
		}
	}
	return order
}

func (boot *CSRFMiddlewareBootstraper) BindMiddleware(server RestServer, config middlewareconfig.MiddlewareConf) {

	CSRFValue, ok := config["csrf"]
	if !ok {
		//用户没有配置
		logger.Info("the CSRF Middleware is not configured, ignore")
		return
	}
	logger.Infof("enable the CSRF Middleware, config is %v", CSRFValue)

	boot.CSRFConfig = &CSRFMiddlewareConfig{}

	if err := transformInterfaceToObject(CSRFValue, boot.CSRFConfig); err != nil {
		logger.Infof("the CSRF Middleware is error config: %s, ignore", err.Error())
		return
	}

	//用户配置启用
	InitProperties()
	server.Use(CSRF())

}

func CSRF() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if skipFilter(c) {
				return next(c)
			}

			tokenFromCookie := tokenFromCookie(c)
			if len(tokenFromCookie) == 0 {
				return next(c)
			}

			tokenFromQuery, err := tokenFromQuery(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden, "csrftoken not found in query params")
			}

			if !isTokenValid(tokenFromCookie, tokenFromQuery) {
				return echo.NewHTTPError(http.StatusForbidden, "invalid csrftoken")
			}

			return next(c)
		}
	}
}

func tokenFromCookie(c echo.Context) string {
	cookieNameWithPort := CSRFTOKEN

	ipPort := c.Request().Header.Get("Gateway-Host")
	if len(ipPort) == 0 {
		ipPort = c.Request().Host
	}
	if len(ipPort) != 0 {
		index := strings.LastIndex(ipPort, ":")
		cookieNameWithPort = cookieNameWithPort + "-" + ipPort[index+1:]
	}

	cookieValueDefault := ""
	cookieValueWithPort := ""
	for _, cookie := range c.Cookies() {
		if strings.Compare(cookie.Name, CSRFTOKEN) == 0 {
			cookieValueDefault = cookie.Value
		}
		if strings.Compare(cookie.Name, cookieNameWithPort) == 0 {
			cookieValueWithPort = cookie.Value
		}
	}

	if len(cookieValueWithPort) > 0 {
		return cookieValueWithPort
	}

	return cookieValueDefault
}

func tokenFromQuery(c echo.Context) (string, error) {
	token := c.Request().Header.Get(CSRFTOKEN)
	if token == "" {
		token = c.QueryParam(CSRFTOKEN)
	}
	if token == "" {
		return "", errors.New("missing csrf token in the query string")
	}
	return token, nil
}

func isTokenValid(t1, t2 string) bool {
	return subtle.ConstantTimeCompare([]byte(t1), []byte(t2)) == 1
}

func skipFilter(c echo.Context) bool {

	path := c.Request().URL.Path
	if isSwagger(path) {
		return true
	}
	if !isExtent(c.Request()) {
		return true
	}
	if isLoginOrLogoutUrl(path) {
		return true
	}
	if isAccessToken(c.Request()) {
		return true
	}
	if matchSpecialUri(path) {
		return true
	}
	return false

}
func isSwagger(uri string) bool {

	if uri != "" && strings.Contains(uri, "swagger") {
		return true
	}
	return false
}

func isExtent(request *http.Request) bool {
	extent := request.Header.Get(Z_EXTENT)
	if extent == "" {
		return false
	}
	return true
}

func isLoginOrLogoutUrl(uri string) bool {
	return strings.Contains(uri, OAUTH_CLIENT_LOGIN) || strings.Contains(uri, OAUTH_CLIENT_LOGOUT)
}

func isAdmin(userName string) bool {
	return userName == USER_ADMIN
}

func isCometd(uri string) bool {
	//TODO Java代码写的是 可能是包含cometd  也有可能是包含并且不以cometd开头的
	//我觉得应该是包含
	return uri != "" && strings.Contains(uri, "cometd")
}

func isAccessToken(request *http.Request) bool {
	opc := request.Form.Get(OAUTH_PASSWORD_CREDENTIAL)
	if opc == "" {
		return false
	}
	isAccessToken, err := strconv.ParseBool(opc)
	if err != nil {
		//TODO
		logger.Info(err)
	}
	return isAccessToken
}

var (
	PropertiesMap         = make(map[string]string)
	NeedFilterSuffixSlice = make([]string, 0, 16)
)

func matchSpecialUri(uri string) bool {

	if uri == "" {
		return false
	}

	if "*.*" == uri {
		return true
	}

	if len(NeedFilterSuffixSlice) == 0 {
		initNeedFilterSuffixSlice()
	}
	for _, suffix := range NeedFilterSuffixSlice {
		if strings.Contains(uri, suffix) {
			return true
		}
	}
	return false
}

func initNeedFilterSuffixSlice() {
	defaultFilterSuffix := strings.Split(DEFAULT_FILTER_SUFFIX, ",")
	for _, value := range defaultFilterSuffix {
		NeedFilterSuffixSlice = append(NeedFilterSuffixSlice, value)
	}
	confFilterSuffixes, ok := PropertiesMap["filterFree"]

	if ok {
		confFilterSuffix := strings.Split(confFilterSuffixes, ",")
		for _, v := range confFilterSuffix {
			NeedFilterSuffixSlice = append(NeedFilterSuffixSlice, v)
		}
	}
}

func InitProperties() {
	InitPropertiesFile(c.GetConfPath(), "auth-config.properties")
}

func InitPropertiesFile(path, fileName string) bool {

	buffer := bytes.Buffer{}
	buffer.WriteString(path)
	buffer.WriteRune(os.PathSeparator)
	buffer.WriteString(fileName)
	filePath := buffer.String()

	f, err := os.Open(filepath.Clean(filePath))
	defer f.Close()

	if err != nil {
		return false
	}

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return false
		}

		s := strings.TrimSpace(string(b))

		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}

		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}

		PropertiesMap[key] = value
	}
	return true
}
