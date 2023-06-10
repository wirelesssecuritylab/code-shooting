package middleware

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	yml "gopkg.in/yaml.v2"
)

func TestRecoverEnable(t *testing.T) {
	Convey("Given a middleware config(recover is configured)", t, func() {
		confContent := `recover:
    name: recover
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the recover middleware config", func() {
			var boot IMiddlewareBootstraper = &RecoverMiddlewareBootstraper{}
			enable := boot.Enable(mwConf)

			Convey("Then enable should be true", func() {
				So(enable, ShouldEqual, true)
			})
		})
	})
}

func TestRecoverOrder(t *testing.T) {
	Convey("Given a middleware config(recover is configured)", t, func() {
		confContent := `recover:
    name: recover
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the recover middleware order", func() {
			var boot IMiddlewareBootstraper = &RecoverMiddlewareBootstraper{}
			order := boot.Order(mwConf)

			Convey("Then order should be 101", func() {
				So(order, ShouldEqual, 101)
			})
		})
	})
}

func TestRecover(t *testing.T) {
	Convey("Given echo and a recover middleware config", t, func() {
		setLogger()
		confContent := `recover:
    name: recover
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			panic("test")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &RecoverMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with handler panic", func() {
			hw := []byte("Hello, World!")
			status, resp := sendRequest(http.MethodPost, "/", hw, e)

			Convey("Then status should be http.StatusInternalServerError and response should contain 'Internal Server Error'", func() {
				So(status, ShouldEqual, http.StatusInternalServerError)
				So(resp, ShouldContainSubstring, "Internal Server Error")
			})
		})
	})

	Convey("Given echo and a recover middleware config(error config)", t, func() {
		setLogger()
		confContent := `recover:
    name: recover
    order: 1
    stack_size:
    - 2B`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			panic("test")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &RecoverMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with handler panic", func() {
			hw := []byte("Hello, World!")
			status, resp := sendRequest(http.MethodPost, "/", hw, e)

			Convey("Then status should be http.StatusInternalServerError and response should contain 'Internal Server Error'", func() {
				So(status, ShouldEqual, http.StatusInternalServerError)
				So(resp, ShouldContainSubstring, "Internal Server Error")
			})
		})
	})

}
