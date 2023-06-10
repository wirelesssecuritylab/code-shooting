package middleware

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/thoas/stats"
	yml "gopkg.in/yaml.v2"
)

func TestStatsEnable(t *testing.T) {
	Convey("Given a middleware config(stats is configured)", t, func() {
		confContent := `stats:
    name: stats
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the stats middleware config", func() {
			var boot IMiddlewareBootstraper = &StatsMiddlewareBootstraper{}
			enable := boot.Enable(mwConf)

			Convey("Then enable should be true", func() {
				So(enable, ShouldEqual, true)
			})
		})
	})
}

func TestStatsOrder(t *testing.T) {
	Convey("Given a middleware config(stats is configured)", t, func() {
		confContent := `stats:
    name: stats
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the stats middleware order", func() {
			var boot IMiddlewareBootstraper = &StatsMiddlewareBootstraper{}
			order := boot.Order(mwConf)

			Convey("Then order should be 101", func() {
				So(order, ShouldEqual, 101)
			})
		})
	})
}

func TestStats(t *testing.T) {
	Convey("Given echo and a stats middleware config", t, func() {
		setLogger()
		confContent := `stats:
    name: stats
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &StatsMiddlewareBootstraper{Middleware: stats.New()}
		boot.BindMiddleware(server, mwConf)

		Convey("When post a request with wildcard origin", func() {
			hw := []byte("Hello, World!")
			status, resp := sendRequest(http.MethodPost, "/", hw, server.server)

			Convey("Then status should be http.StatusNotFound and response should be 'Hello, World!'", func() {
				So(status, ShouldEqual, http.StatusOK)
				So(resp, ShouldEqual, "Hello, World!")
			})
		})
	})

}
