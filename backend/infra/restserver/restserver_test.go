package restserver

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"code-shooting/infra/restserver/internal"

	. "github.com/agiledragon/gomonkey"
	echo "github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	yml "gopkg.in/yaml.v2"

	"code-shooting/infra/restserver/middleware"
)

func getRestServerConf() *internal.RestServerConf {
	return &internal.RestServerConf{
		Name:     "restserver",
		RootPath: "root",
		HttpServer: internal.HttpServerConf{
			Protocol:          "",
			Addr:              "127.0.0.1:0",
			ReadTimeout:       180 * time.Second,
			ReadHeaderTimeout: 0,
			WriteTimeout:      60 * time.Second,
			IdleTimeout:       0,
			MaxHeaderBytes:    16384,
			CertFile:          "",
			KeyFile:           "",
		},
		Listener: internal.ListenerConf{
			Addr:           "127.0.0.1:0",
			MaxConnections: 0,
		},
		Middlewares: make(map[string]interface{}),
	}
}

func TestRestServer(t *testing.T) {

	Convey("Given a rest-server ", t, func() {
		patches := ApplyFunc(buildHttpServer, func(_ *internal.HttpServerConf) (*http.Server, error) {
			return nil, nil
		})
		defer patches.Reset()

		e := echo.New()
		restServer, _ := newRestServer(e, getRestServerConf())
		restServer.RootGroupBox = &RootGroupBox{RootGroup: e.Group("/root"), RootPath: "/root"}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := restServer.NewContext(req, rec)
		Convey("When get the router", func() {

			r := restServer.Router()
			// DefaultHTTPErrorHandler
			restServer.DefaultHTTPErrorHandler(errors.New("error"), c)
			Convey("Then router is not nil ", func() {
				So(r, ShouldNotBeNil)
				So(rec.Code, ShouldEqual, http.StatusInternalServerError)
			})

		})
	})
}

func request(method, path string, e *echo.Echo) (int, string) {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func TestRestServerRoute(t *testing.T) {
	Convey("Given a rest-server ", t, func() {
		patches := ApplyFunc(buildHttpServer, func(_ *internal.HttpServerConf) (*http.Server, error) {
			return nil, nil
		})
		defer patches.Reset()

		e := echo.New()
		restServer, _ := newRestServer(e, getRestServerConf())
		restServer.RootGroupBox = &RootGroupBox{RootGroup: e.Group("/root"), RootPath: "/root"}

		Convey("When create a group of urls", func() {
			h := func(Context) error { return nil }
			rcon := restServer.CONNECT("/", h)
			rdel := restServer.DELETE("/", h)
			rget := restServer.GET("/", h)
			rhead := restServer.HEAD("/", h)
			ropt := restServer.OPTIONS("/", h)
			rpatch := restServer.PATCH("/", h)
			rpost := restServer.POST("/", h)
			rput := restServer.PUT("/", h)
			rtrace := restServer.TRACE("/", h)
			rstat := restServer.Static("/static", "/tmp")
			rfile := restServer.File("/walle", "_fixture/images//walle.png")
			radd := restServer.Add(http.MethodGet, "/add", h)

			Convey("Then get routers should ", func() {
				routes := restServer.Routes()
				So(routes, ShouldContain, rcon)
				So(routes, ShouldContain, rdel)
				So(routes, ShouldContain, rget)
				So(routes, ShouldContain, rhead)
				So(routes, ShouldContain, ropt)
				So(routes, ShouldContain, rpatch)
				So(routes, ShouldContain, rpost)
				So(routes, ShouldContain, rput)
				So(routes, ShouldContain, rtrace)
				So(routes, ShouldContain, rstat)
				So(routes, ShouldContain, rfile)
				So(routes, ShouldContain, radd)
			})
		})

		Convey("When register route for multiple HTTP methods", func() {
			h := func(Context) error { return nil }
			mResp200 := func(next HandlerFunc) HandlerFunc {
				return func(c Context) error {
					return c.NoContent(200)
				}
			}
			mResp400 := func(next HandlerFunc) HandlerFunc {
				return func(c Context) error {
					return c.NoContent(400)
				}
			}
			restServer.Any("/any", h, mResp200)
			restServer.Match([]string{http.MethodGet, http.MethodPost}, "/match", h, mResp400)

			Convey("Then response status should be expect", func() {
				c, _ := request(http.MethodGet, "/any", e)
				So(c, ShouldEqual, 200)
				c, _ = request(http.MethodPost, "/match", e)
				So(c, ShouldEqual, 400)
			})
		})

		Convey("When create a url group", func() {
			h := func(Context) error { return nil }
			mResp204 := func(next HandlerFunc) HandlerFunc {
				return func(c Context) error {
					return c.NoContent(204)
				}
			}
			rgroup := restServer.Group("/group")
			rgroup.GET("/204", h, mResp204)

			Convey("Then response status should be expect", func() {
				c, _ := request(http.MethodGet, "/group/204", e)
				So(c, ShouldEqual, 204)
			})
		})

		Convey("When register a route ", func() {
			file := func(Context) error { return nil }
			registeredUrlName := "/static/file"
			restServer.GET(registeredUrlName, file)

			registeredUriName := "/static/photo"
			photo := func(Context) error { return nil }
			restServer.GET(registeredUriName, photo)

			foobar := func(Context) error { return nil }
			reverseName := "/users/:id"
			restServer.GET(reverseName, foobar).Name = "foobar"

			Convey("Then get the handler name should be registered function name", func() {
				urlName := restServer.URL(file)
				So(urlName, ShouldEqual, registeredUrlName)

				uriName := restServer.URI(photo)
				So(uriName, ShouldEqual, registeredUriName)

				reversename := restServer.Reverse("foobar", 1)
				So(reversename, ShouldEqual, "/users/1")

			})
		})

	})
}

func waitForServerStart(e *echo.Echo, errChan <-chan error, isTLS bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			var addr net.Addr
			if isTLS && e.TLSListener != nil {
				addr = e.TLSListener.Addr()
			} else if e.Listener != nil {
				addr = e.Listener.Addr()
			}
			if addr != nil && strings.Contains(addr.String(), ":") {
				return nil // was started
			}
		case err := <-errChan:
			if err == http.ErrServerClosed {
				return nil
			}
			return err
		}
	}
}

func TestStartHTTPServer(t *testing.T) {

	Convey("Given a rest-server ", t, func() {
		patches := ApplyFunc(buildHttpServer, func(_ *internal.HttpServerConf) (*http.Server, error) {
			return nil, nil
		})
		defer patches.Reset()

		e := echo.New()
		restServer, _ := newRestServer(e, getRestServerConf())
		restServer.RootGroupBox = &RootGroupBox{RootGroup: e.Group("/root"), RootPath: "/root"}
		errChan := make(chan error)

		Convey("When start an HTTP Server with Start", func() {
			go func() {
				err := restServer.Start("127.0.0.1:8081")
				if err != nil {
					errChan <- err
				}
			}()

			Convey("Then start server success", func() {
				err := waitForServerStart(restServer.server, errChan, true)
				So(err, ShouldBeNil)

				err = restServer.Close()
				So(err, ShouldBeNil)

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				err = restServer.Shutdown(ctx)
				So(err, ShouldBeNil)
			})
		})

		Convey("When start an HTTP Server with StartTLS without certFile and keyFile", func() {

			go func() {
				err := restServer.StartTLS("127.0.0.1:8081", "not existing", "not existing")
				if err != nil {
					errChan <- err
				}
			}()

			Convey("Then start server failed", func() {
				err := waitForServerStart(restServer.server, errChan, true)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "no such file or directory")
			})
		})

		Convey("When start an HTTP Server with StartAutoTLS without address", func() {
			errChan := make(chan error)
			go func() {
				errChan <- restServer.StartAutoTLS("nope")
			}()

			Convey("Then start server failed", func() {
				err := waitForServerStart(restServer.server, errChan, true)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "missing port in address")
			})
		})

		Convey("When start an HTTP Server with StartServer without address", func() {
			errChan := make(chan error)
			server := new(http.Server)
			server.Addr = "nope"
			go func() {
				errChan <- restServer.StartServer(server)
			}()

			Convey("Then start server failed", func() {
				err := waitForServerStart(restServer.server, errChan, true)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "missing port in address")
			})
		})
	})
}

func TestServeHTTP(t *testing.T) {

	Convey("Given a rest-server ", t, func() {
		patches := ApplyFunc(buildHttpServer, func(_ *internal.HttpServerConf) (*http.Server, error) {
			return nil, nil
		})
		defer patches.Reset()

		e := echo.New()
		restServer, _ := newRestServer(e, getRestServerConf())
		restServer.RootGroupBox = &RootGroupBox{RootGroup: e.Group("/root"), RootPath: "/root"}

		Convey("When serve http with a not existed url", func() {
			req := httptest.NewRequest(http.MethodGet, "/root", nil)
			rec := httptest.NewRecorder()
			restServer.ServeHTTP(rec, req)
			Convey("Then rec.Code should be 404", func() {
				So(rec.Code, ShouldEqual, 404)
			})

		})

	})
}
func TestContext(t *testing.T) {

	Convey("Given a rest-server ", t, func() {
		patches := ApplyFunc(buildHttpServer, func(_ *internal.HttpServerConf) (*http.Server, error) {
			return nil, nil
		})
		defer patches.Reset()

		e := echo.New()
		restServer, _ := newRestServer(e, getRestServerConf())
		restServer.RootGroupBox = &RootGroupBox{RootGroup: e.Group("/root"), RootPath: "/root"}

		Convey("When get an empty Context instance", func() {

			c := restServer.AcquireContext()

			Convey("Then type of the instance is context", func() {
				So(c, ShouldHaveSameTypeAs, restServer.NewContext(nil, nil))

				restServer.ReleaseContext(c)
			})
		})
	})
}

func TestMiddleware(t *testing.T) {

	Convey("Given a rest-server ", t, func() {
		patches := ApplyFunc(buildHttpServer, func(_ *internal.HttpServerConf) (*http.Server, error) {
			return nil, nil
		})
		defer patches.Reset()

		e := echo.New()
		restServer, _ := newRestServer(e, getRestServerConf())
		restServer.RootGroupBox = &RootGroupBox{RootGroup: e.Group("/root"), RootPath: "/root"}

		Convey("When register a route with a middleware which is run before router", func() {

			buf := new(bytes.Buffer)
			restServer.Pre(func(next HandlerFunc) HandlerFunc {
				return func(c Context) error {
					buf.WriteString("-1")
					return next(c)
				}
			})

			restServer.GET("/", func(c Context) error {
				return c.String(http.StatusOK, "OK")
			})

			c, b := request(http.MethodGet, "/", restServer.server)

			Convey("Then status should be ok and response should be ok", func() {
				So(buf.String(), ShouldEqual, "-1")
				So(c, ShouldEqual, http.StatusOK)
				So(b, ShouldEqual, "OK")

			})
		})

		Convey("When register a route with a middleware which is run after router", func() {

			restServer.Use(func(next HandlerFunc) HandlerFunc {
				return func(c Context) error {
					return errors.New("error")
				}
			})
			e.GET("/", func(c Context) error {
				return echo.ErrNotFound
			})
			c, _ := request(http.MethodGet, "/", restServer.server)

			Convey("Then response should be StatusInternalServerError", func() {

				So(c, ShouldEqual, http.StatusInternalServerError)

			})
		})
	})

}

func sendRequest(method, path string, body []byte, e *echo.Echo) (int, string) {

	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func TestRateLimitMiddleware(t *testing.T) {

	Convey("Given a rest-server ", t, func() {
		patches := ApplyFunc(buildHttpServer, func(_ *internal.HttpServerConf) (*http.Server, error) {
			return nil, nil
		})
		defer patches.Reset()

		e := echo.New()
		e.POST("/hello", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})
		restServer, _ := newRestServer(e, getRestServerConf())
		restServer.RootGroupBox = &RootGroupBox{RootGroup: e.Group("/root"), RootPath: "/root"}

		Convey("When bind rateLimitMiddleware", func() {
			confContent := `ratelimit:
    name: ratelimit
    maxrequests: 2
    requestspersec: 1
    order: 1`

			mwConf := make(map[string]interface{})
			yml.Unmarshal([]byte(confContent), mwConf)

			var boot middleware.IMiddlewareBootstraper = &middleware.RateLimitMiddlewareBootstraper{}
			boot.BindMiddleware(restServer, mwConf)

			Convey("When post 7 requests in 1s\n", func() {
				var refusedNum, acceptedNum int32
				hw := []byte("Hello, World!")
				wg := sync.WaitGroup{}
				for i := 0; i < 7; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						status, _ := sendRequest(http.MethodPost, "/hello", hw, restServer.server)
						if status == http.StatusTooManyRequests {
							atomic.AddInt32(&refusedNum, 1)
						} else if status == http.StatusOK {
							atomic.AddInt32(&acceptedNum, 1)
						}
					}()
				}

				wg.Wait()

				Convey("Then 2 requests should be accepted and 5 requests should be refused\n", func() {
					So(acceptedNum, ShouldEqual, 2)
					So(refusedNum, ShouldEqual, 5)
				})
			})
		})
	})
}
