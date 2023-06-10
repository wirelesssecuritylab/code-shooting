package middleware

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	yml "gopkg.in/yaml.v2"
)

func TestBodyLimitEnable(t *testing.T) {
	Convey("Given a middleware config(bodylimit is configured)", t, func() {
		setLogger()
		confContent := `bodylimit:
    name: bodylimit
    limit: 2M
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the bodylimit middleware config", func() {
			var boot IMiddlewareBootstraper = &BodyLimitMiddlewareBootstraper{}
			enable := boot.Enable(mwConf)

			Convey("Then enable should be true", func() {
				So(enable, ShouldEqual, true)
			})
		})
	})
}

func TestBodyLimitOrder(t *testing.T) {
	Convey("Given a middleware config(bodylimit is configured)", t, func() {
		confContent := `bodylimit:
    name: bodylimit
    limit: 2M
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the bodylimit middleware order", func() {
			var boot IMiddlewareBootstraper = &BodyLimitMiddlewareBootstraper{}
			order := boot.Order(mwConf)

			Convey("Then order should be 101", func() {
				So(order, ShouldEqual, 101)
			})
		})
	})
}

func TestBodyLimit(t *testing.T) {
	Convey("Given echo and a bodylimit middleware config(limit is big)", t, func() {
		confContent := `bodylimit:
    name: bodylimit
    limit: 2M
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/hello", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()

		var boot IMiddlewareBootstraper = &BodyLimitMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with body overlimit", func() {
			hw := []byte("Hello, World!")
			status, resp := sendRequest(http.MethodPost, "/hello", hw, e)

			Convey("Then status should be http.StatusOK and response should be 'Hello, World!'", func() {
				So(status, ShouldEqual, http.StatusOK)
				So(string(resp), ShouldEqual, "Hello, World!")
			})
		})
	})

	Convey("Given echo and a bodylimit middleware config(limit is small)", t, func() {
		confContent := `bodylimit:
    name: bodylimit
    limit: 2B
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/hello", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &BodyLimitMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with body overlimit", func() {
			hw := []byte("Hello, World!")
			status, resp := sendRequest(http.MethodPost, "/hello", hw, e)

			Convey("Then status should be http.StatusRequestEntityTooLarge and response should contain 'Request Entity Too Large'", func() {
				So(status, ShouldEqual, http.StatusRequestEntityTooLarge)
				So(string(resp), ShouldContainSubstring, "Request Entity Too Large")
			})
		})
	})

}
