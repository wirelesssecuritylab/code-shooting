package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	yml "gopkg.in/yaml.v2"
)

func TestRequestIDEnable(t *testing.T) {
	Convey("Given a middleware config(requestid is configured)", t, func() {
		confContent := `requestid:
    name: requestid
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the requestid middleware config", func() {
			var boot IMiddlewareBootstraper = &RequestIDMiddlewareBootstraper{}
			enable := boot.Enable(mwConf)

			Convey("Then enable should be true", func() {
				So(enable, ShouldEqual, true)
			})
		})
	})
}

func TestRequestIDOrder(t *testing.T) {
	Convey("Given a middleware config(requestid is configured)", t, func() {
		confContent := `requestid:
    name: requestid
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the requestid middleware order", func() {
			var boot IMiddlewareBootstraper = &RequestIDMiddlewareBootstraper{}
			order := boot.Order(mwConf)

			Convey("Then order should be 101", func() {
				So(order, ShouldEqual, 101)
			})
		})
	})
}

func TestRequestID(t *testing.T) {
	Convey("Given echo and a requestid middleware config", t, func() {
		setLogger()
		confContent := `requestid:
    name: requestid
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &RequestIDMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request", func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			Convey("Then requestid len should be 32 and status should be http.StatusOK and response should be 'Hello, World!'", func() {
				So(len(rec.Header().Get(echo.HeaderXRequestID)), ShouldEqual, 32)
				So(rec.Code, ShouldEqual, http.StatusOK)
				So(rec.Body.String(), ShouldEqual, "Hello, World!")
			})
		})
	})

	Convey("Given echo and a requestid middleware config(error config", t, func() {
		setLogger()
		confContent := `requestid:
    name: requestid
    order: 1
    generator: none`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &RequestIDMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request", func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			Convey("Then use default generator, requestid len should be 32 and status should be http.StatusOK and response should be 'Hello, World!'", func() {
				So(len(rec.Header().Get(echo.HeaderXRequestID)), ShouldEqual, 32)
				So(rec.Code, ShouldEqual, http.StatusOK)
				So(rec.Body.String(), ShouldEqual, "Hello, World!")
			})
		})
	})

}
