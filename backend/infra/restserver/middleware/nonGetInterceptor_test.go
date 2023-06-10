package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	yml "gopkg.in/yaml.v2"
)

func TestNonGetInterceptorEnable(t *testing.T) {
	Convey("Given a middleware config(nonGetInterceptor is configured)", t, func() {
		confContent := `nonGetInterceptor:
    name: nonGetInterceptor
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the nonGetInterceptor middleware config", func() {
			var boot IMiddlewareBootstraper = &NonGetInterceptorBootstraper{}
			enable := boot.Enable(mwConf)

			Convey("Then enable should be true", func() {
				So(enable, ShouldEqual, true)
			})
		})
	})
}

func TestNonGetInterceptorOrder(t *testing.T) {
	Convey("Given a middleware config(nonGetInterceptor is configured)", t, func() {
		confContent := `nonGetInterceptor:
    name: nonGetInterceptor
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the nonGetInterceptor middleware order", func() {
			var boot IMiddlewareBootstraper = &NonGetInterceptorBootstraper{}
			order := boot.Order(mwConf)

			Convey("Then order should be 101", func() {
				So(order, ShouldEqual, 101)
			})
		})
	})
}

func TestNonGetInterceptor(t *testing.T) {
	Convey("Given echo and a nonGetInterceptor middleware config of master", t, func() {
		setLogger()
		confContent := `nonGetInterceptor:
    name: nonGetInterceptor
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &NonGetInterceptorBootstraper{IsMaster: true}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request blocked", func() {
			hw := []byte("Hello, World!")
			status, resp := sendRequest(http.MethodPost, "/", hw, e)

			Convey("Then status should be http.StatusOK and response should be 'Hello, World!'", func() {
				So(status, ShouldEqual, http.StatusOK)
				So(resp, ShouldEqual, "Hello, World!")
			})
		})
	})

	Convey("Given echo and a nonGetInterceptor middleware config with allowlist", t, func() {
		setLogger()
		confContent := `nonGetInterceptor:
    name: nonGetInterceptor
    allowlist:
    - /allow/user
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/allow/user", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &NonGetInterceptorBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request allowed", func() {
			hw := []byte("Hello, World!")
			status, resp := sendRequest(http.MethodPost, "/allow/user", hw, e)

			Convey("Then status should be http.StatusOK and response should be 'Hello, World!'", func() {
				So(status, ShouldEqual, http.StatusOK)
				So(resp, ShouldEqual, "Hello, World!")
			})
		})
	})

	Convey("Given echo and a nonGetInterceptor middleware config without allowlist", t, func() {
		setLogger()
		confContent := `nonGetInterceptor:
    name: nonGetInterceptor
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.GET("/allow/user", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &NonGetInterceptorBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with get method", func() {
			hw := []byte("Hello, World!")
			status, resp := sendRequest(http.MethodGet, "/allow/user", hw, e)

			Convey("Then status should be http.StatusOK and response should be 'Hello, World!'", func() {
				So(status, ShouldEqual, http.StatusOK)
				So(resp, ShouldEqual, "Hello, World!")
			})
		})
	})

	Convey("Given echo and a nonGetInterceptor middleware config without allowlist", t, func() {
		setLogger()
		confContent := `nonGetInterceptor:
    name: nonGetInterceptor
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/allow/user", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &NonGetInterceptorBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with language-option en-us", func() {

			req := httptest.NewRequest(http.MethodPost, "/allow/user", bytes.NewReader([]byte("Hello, World!")))
			rec := httptest.NewRecorder()
			req.Header.Set("language-option", "en-us")
			e.ServeHTTP(rec, req)

			Convey("Then status should be http.StatusOK and response should contain 'not allowed'", func() {
				So(rec.Code, ShouldEqual, http.StatusForbidden)
				So(rec.Body.String(), ShouldContainSubstring, "not allowed")
			})
		})

		Convey("When post a request with language-option zh-cn", func() {

			req := httptest.NewRequest(http.MethodPost, "/allow/user", bytes.NewReader([]byte("Hello, World!")))
			rec := httptest.NewRecorder()
			req.Header.Set("language-option", "zh-cn")
			e.ServeHTTP(rec, req)

			Convey("Then status should be http.StatusForbidden and response should contain '不允许'", func() {
				So(rec.Code, ShouldEqual, http.StatusForbidden)
				So(rec.Body.String(), ShouldContainSubstring, "不允许")
			})
		})
	})

	Convey("Given echo and a nonGetInterceptor middleware config with error allowlist", t, func() {
		setLogger()
		confContent := `nonGetInterceptor:
    name: nonGetInterceptor
    allowlist: /allow/user
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/allow/user", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &NonGetInterceptorBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request allowed", func() {
			hw := []byte("Hello, World!")
			status, resp := sendRequest(http.MethodPost, "/allow/user", hw, e)

			Convey("Then status should be http.StatusOK and response should be 'Hello, World!'", func() {
				So(status, ShouldEqual, http.StatusOK)
				So(resp, ShouldEqual, "Hello, World!")
			})
		})
	})
}
