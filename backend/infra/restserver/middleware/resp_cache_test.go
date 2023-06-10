package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	yml "gopkg.in/yaml.v2"
)

func TestRespCacheEnable(t *testing.T) {
	Convey("Given a middleware config(respcache is configured)", t, func() {
		confContent := `respcache:
    name: respcache
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the respcache middleware config", func() {
			var boot IMiddlewareBootstraper = &RespCacheMiddlewareBootstraper{}
			enable := boot.Enable(mwConf)

			Convey("Then enable should be true", func() {
				So(enable, ShouldEqual, true)
			})
		})
	})
}

func TestRespCacheOrder(t *testing.T) {
	Convey("Given a middleware config(respcache is configured)", t, func() {
		confContent := `respcache:
    name: respcache
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the respcache middleware order", func() {
			var boot IMiddlewareBootstraper = &RespCacheMiddlewareBootstraper{}
			order := boot.Order(mwConf)

			Convey("Then order should be 101", func() {
				So(order, ShouldEqual, 101)
			})
		})
	})
}

func TestRespCache(t *testing.T) {
	Convey("Given echo and a respcache middleware config(default config)", t, func() {
		setLogger()
		confContent := `respcache:
    name: respcache
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &RespCacheMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request ", func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			Convey("Then status should be http.StatusOK and response should be 'Hello, World!'", func() {
				So(rec.Code, ShouldEqual, http.StatusOK)
				So(rec.Body.String(), ShouldEqual, "Hello, World!")
			})
		})
	})

	Convey("Given echo and a respcache middleware config(user config)", t, func() {
		setLogger()
		confContent := `respcache:
    name: respcache
    order: 1
    cachecontrol: only-if-cached
    expires: 1
    pragma: public
    xcontenttypeoptions: text/html;charset=utf-8`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &RespCacheMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request ", func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			Convey("Then status should be http.StatusOK and response should be 'Hello, World!'", func() {
				So(rec.Code, ShouldEqual, http.StatusOK)
				So(rec.Body.String(), ShouldEqual, "Hello, World!")
			})
		})
	})

	Convey("Given echo and a respcache middleware config(error config)", t, func() {
		setLogger()
		confContent := `respcache:
    name: 
    - respcache
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &RespCacheMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request ", func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			Convey("Then status should be http.StatusOK and response should be 'Hello, World!'", func() {
				So(rec.Code, ShouldEqual, http.StatusOK)
				So(rec.Body.String(), ShouldEqual, "Hello, World!")
			})
		})
	})

}
