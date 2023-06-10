package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	yml "gopkg.in/yaml.v2"
)

func TestHeaderEnable(t *testing.T) {
	Convey("Given a middleware config(header is configured)", t, func() {
		confContent := `header:
    name: header
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the header middleware config", func() {
			var boot IMiddlewareBootstraper = &HeaderMiddlewareBootstraper{}
			enable := boot.Enable(mwConf)

			Convey("Then enable should be true", func() {
				So(enable, ShouldEqual, true)
			})
		})
	})
}

func TestHeaderOrder(t *testing.T) {
	Convey("Given a middleware config(header is configured)", t, func() {
		confContent := `header:
    name: header
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the header middleware order", func() {
			var boot IMiddlewareBootstraper = &HeaderMiddlewareBootstraper{}
			order := boot.Order(mwConf)

			Convey("Then order should be 101", func() {
				So(order, ShouldEqual, 101)
			})
		})
	})
}

func TestHeader(t *testing.T) {
	Convey("Given echo and a header middleware config", t, func() {
		setLogger()
		confContent := `header:
    name: header
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &HeaderMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with traceid", func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			req.Header.Set(TRACE_ID_NAME, fmt.Sprintf("%06v", time.Now().UnixNano()))
			e.ServeHTTP(rec, req)

			Convey("Then status should be http.StatusNotFound and response should be 'Hello, World!'", func() {
				So(rec.Code, ShouldEqual, http.StatusOK)
				So(rec.Body.String(), ShouldEqual, "Hello, World!")
			})
		})
	})

}
