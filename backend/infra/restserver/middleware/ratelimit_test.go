package middleware

import (
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	yml "gopkg.in/yaml.v2"
)

func TestRateLimitEnable(t *testing.T) {
	Convey("Given a middleware config(ratelimit is configured)", t, func() {
		setLogger()
		confContent := `ratelimit:
    name: ratelimit
    maxrequests: 2
    requestspersec: 1
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the ratelimit middleware config", func() {
			var boot IMiddlewareBootstraper = &RateLimitMiddlewareBootstraper{}
			enable := boot.Enable(mwConf)

			Convey("Then enable should be true", func() {
				So(enable, ShouldEqual, true)
			})
		})
	})
}

func TestRateLimitOrder(t *testing.T) {
	Convey("Given a middleware config(ratelimit is configured)", t, func() {
		confContent := `ratelimit:
    name: ratelimit
    maxrequests: 2
    requestspersec: 1
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		Convey("When get the ratelimit middleware order", func() {
			var boot IMiddlewareBootstraper = &RateLimitMiddlewareBootstraper{}
			order := boot.Order(mwConf)

			Convey("Then order should be 101", func() {
				So(order, ShouldEqual, 101)
			})
		})
	})
}

func TestRateLimit(t *testing.T) {
	Convey("Given echo and a ratelimit middleware with requestspersec is 3\n", t, func() {
		confContent := `ratelimit:
    name: ratelimit
    maxrequests: 3
    requestspersec: 3
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/hello", func(c echo.Context) error {
			time.Sleep(1 * time.Second)
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &RateLimitMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post 7 requests in 1s\n", func() {
			var refusedNum, acceptedNum int32
			hw := []byte("Hello, World!")
			wg := sync.WaitGroup{}
			for i := 0; i < 7; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					status, _ := sendRequest(http.MethodPost, "/hello", hw, e)
					if status == http.StatusTooManyRequests {
						atomic.AddInt32(&refusedNum, 1)
					} else if status == http.StatusOK {
						atomic.AddInt32(&acceptedNum, 1)
					}
				}()
			}
			wg.Wait()

			Convey("Then 3 requests should be accepted and 4 requests should be refused\n", func() {
				So(acceptedNum, ShouldEqual, 3)
				So(refusedNum, ShouldEqual, 4)
			})
		})
	})

	Convey("Given echo and a ratelimit middleware with requestspersec is 3\n", t, func() {
		confContent := `ratelimit:
    name: ratelimit
    maxrequests: 3
    requestspersec: 3
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/hello", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &RateLimitMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post 7 requests with 1/s rate\n", func() {
			var refusedNum, acceptedNum int32
			hw := []byte("Hello, World!")
			wg := sync.WaitGroup{}
			for i := 0; i < 7; i++ {
				wg.Add(1)
				time.Sleep(1 * time.Second)
				go func() {
					defer wg.Done()
					status, _ := sendRequest(http.MethodPost, "/hello", hw, e)
					if status == http.StatusTooManyRequests {
						atomic.AddInt32(&refusedNum, 1)
					} else if status == http.StatusOK {
						atomic.AddInt32(&acceptedNum, 1)
					}
				}()
			}
			wg.Wait()

			Convey("Then all 7 requests should be accepted\n", func() {
				So(acceptedNum, ShouldEqual, 7)
				So(refusedNum, ShouldEqual, 0)
			})
		})
	})
}

func TestRateLimitUpdate(t *testing.T) {
	Convey("Given echo and a ratelimit middleware with requestspersec is 3\n", t, func() {
		confContent := `ratelimit:
    name: ratelimit
    maxrequests: 2
    requestspersec: 1
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/hello", func(c echo.Context) error {
			time.Sleep(1 * time.Second)
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &RateLimitMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When ratelimit update with requestspersec is 2, and post 7 requests in 1s\n", func() {
			So(server.SetRateLimitConfig(4, 3), ShouldBeNil)

			time.Sleep(1 * time.Second)

			var refusedNum, acceptedNum int32
			hw := []byte("Hello, World!")
			wg := sync.WaitGroup{}
			for i := 0; i < 7; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					status, _ := sendRequest(http.MethodPost, "/hello", hw, server.server)
					if status == http.StatusTooManyRequests {
						atomic.AddInt32(&refusedNum, 1)
					} else if status == http.StatusOK {
						atomic.AddInt32(&acceptedNum, 1)
					}
				}()
			}
			wg.Wait()

			Convey("Then 4 requests should be accepted and 3 requests should be refused\n", func() {
				So(acceptedNum, ShouldEqual, 4)
				So(refusedNum, ShouldEqual, 3)
			})
		})
	})

	Convey("Given echo and a ratelimit middleware with requestspersec is 1, maxrequests is 3\n", t, func() {
		confContent := `ratelimit:
    name: ratelimit
    maxrequests: 3
    requestspersec: 1
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/hello", func(c echo.Context) error {
			time.Sleep(1 * time.Second)
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)
		defer server.server.Close()
		var boot IMiddlewareBootstraper = &RateLimitMiddlewareBootstraper{}
		boot.BindMiddleware(server, mwConf)

		Convey("When post 9 requests with 3post/1sec \n", func() {
			var refusedNum, acceptedNum int32
			hw := []byte("Hello, World!")
			wg := sync.WaitGroup{}
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						status, _ := sendRequest(http.MethodPost, "/hello", hw, server.server)
						if status == http.StatusTooManyRequests {
							atomic.AddInt32(&refusedNum, 1)
						} else if status == http.StatusOK {
							atomic.AddInt32(&acceptedNum, 1)
						}
					}()
				}
				time.Sleep(1 * time.Second)
			}

			wg.Wait()

			Convey("Then 5 requests should be accepted and 4 requests should be refused\n", func() {
				So(acceptedNum, ShouldEqual, 5)
				So(refusedNum, ShouldEqual, 4)
			})
		})
	})
}
