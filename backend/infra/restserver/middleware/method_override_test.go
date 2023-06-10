package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	yml "gopkg.in/yaml.v2"
)

func TestMethodOverrideEnable(t *testing.T) {
	Convey("Given a middleware config(methodOverride is configured)", t, func() {
		confContent := `methodOverride:
    name: methodOverride
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the methodOverride middleware config", func() {
			var boot IMiddlewareBootstraper = &MethodOverrideMiddlewareBootstraper{}
			enable := boot.Enable(mwConf)

			Convey("Then enable should be true", func() {
				So(enable, ShouldEqual, true)
			})
		})
	})
}

func TestMethodOverrideOrder(t *testing.T) {
	Convey("Given a middleware config(methodOverride is configured)", t, func() {
		confContent := `methodOverride:
    name: methodOverride
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the methodOverride middleware order", func() {
			var boot IMiddlewareBootstraper = &MethodOverrideMiddlewareBootstraper{}
			order := boot.Order(mwConf)

			Convey("Then order should be 101", func() {
				So(order, ShouldEqual, 101)
			})
		})
	})
}

func TestMethodOverrideHeader(t *testing.T) {
	Convey("Given echo and a methodOverride middleware config", t, func() {
		confContent := `methodOverride:
    name: methodOverride
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &MethodOverrideMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with X-HTTP-Method", func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			req.Header.Set(X_HTTP_METHOD, "POST")
			e.ServeHTTP(rec, req)

			Convey("Then status should be http.StatusMethodNotAllowed", func() {
				So(rec.Code, ShouldEqual, http.StatusMethodNotAllowed)
			})
		})
	})

}

func TestMethodOverrideUrl(t *testing.T) {
	Convey("Given echo and a header-method middleware config", t, func() {
		confContent := `methodOverride:
    name: methodOverride
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &MethodOverrideMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with url contains _method", func() {
			req := httptest.NewRequest(http.MethodPost, "/?_method=POST", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			Convey("Then status should be http.StatusMethodNotAllowed", func() {
				So(rec.Code, ShouldEqual, http.StatusMethodNotAllowed)
			})
		})
	})
}
