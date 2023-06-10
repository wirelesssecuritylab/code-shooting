package middleware

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	yml "gopkg.in/yaml.v2"
)

func TestCORSEnable(t *testing.T) {
	Convey("Given a middleware config(cors is configured)", t, func() {
		confContent := `cors:
    name: cors
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the cors middleware config", func() {
			var boot IMiddlewareBootstraper = &CORSMiddlewareBootstraper{}
			enable := boot.Enable(mwConf)

			Convey("Then enable should be true", func() {
				So(enable, ShouldEqual, true)
			})
		})
	})
}

func TestCORSOrder(t *testing.T) {
	Convey("Given a middleware config(cors is configured)", t, func() {
		confContent := `cors:
    name: cors
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the cors middleware order", func() {
			var boot IMiddlewareBootstraper = &CORSMiddlewareBootstraper{}
			order := boot.Order(mwConf)

			Convey("Then order should be 101", func() {
				So(order, ShouldEqual, 101)
			})
		})
	})
}

func TestCORS(t *testing.T) {
	Convey("Given echo and a cors middleware config(error config)", t, func() {
		setLogger()
		confContent := `cors:
    name: cors
    allow_origins: localhost
    allow_headers: Content-Type, Origin
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &CORSMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with wildcard origin", func() {
			status, resp := sendRequest(http.MethodGet, "*", nil, e)

			Convey("Then status should be http.StatusNotFound and response should contain 'Not Found'", func() {
				So(status, ShouldEqual, http.StatusNotFound)
				So(string(resp), ShouldContainSubstring, "Not Found")
			})
		})
	})

	Convey("Given echo and a cors middleware config", t, func() {
		confContent := `cors:
    name: cors
    allow_origins:
    - localhost
    allow_methods:
    - OPTIONS
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &CORSMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with option method", func() {
			status, resp := sendRequest(http.MethodOptions, "/", nil, e)

			Convey("Then status should be http.StatusNotFound and response should contain 'Not Found'", func() {
				So(status, ShouldEqual, http.StatusNoContent)
				So(string(resp), ShouldEqual, "")
			})
		})
	})
}
