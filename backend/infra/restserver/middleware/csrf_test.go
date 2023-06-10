package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	yml "gopkg.in/yaml.v2"
)

func TestCSRFEnable(t *testing.T) {
	Convey("Given a middleware config(csrf is configured)", t, func() {
		confContent := `csrf:
    name: csrf
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the csrf middleware config", func() {
			var boot IMiddlewareBootstraper = &CSRFMiddlewareBootstraper{}
			enable := boot.Enable(mwConf)

			Convey("Then enable should be true", func() {
				So(enable, ShouldEqual, true)
			})
		})
	})
}

func TestCSRFOrder(t *testing.T) {
	Convey("Given a middleware config(csrf is configured)", t, func() {
		confContent := `csrf:
    name: csrf
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the csrf middleware order", func() {
			var boot IMiddlewareBootstraper = &CSRFMiddlewareBootstraper{}
			order := boot.Order(mwConf)

			Convey("Then order should be 101", func() {
				So(order, ShouldEqual, 101)
			})
		})
	})
}

func TestCSRF(t *testing.T) {
	Convey("Given echo and a csrf middleware config", t, func() {
		setLogger()
		confContent := `csrf:
    name: csrf
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &CSRFMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with CSRF token", func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderSetCookie, "_csrf")
			req.Header.Set(Z_EXTENT, "true")
			e.ServeHTTP(rec, req)

			Convey("Then status should be http.StatusOK and response should be 'Hello, World!'", func() {
				So(rec.Code, ShouldEqual, http.StatusOK)
				So(rec.Body.String(), ShouldEqual, "Hello, World!")
			})
		})
	})

}
