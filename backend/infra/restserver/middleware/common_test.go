package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
	yml "gopkg.in/yaml.v2"

	marsconfig "code-shooting/infra/config"
	log "code-shooting/infra/logger"
	"code-shooting/infra/x/test"
)

func setLogger() {
	content := `code-shooting:
    log:
      level: info
      encoder: plain
      outputPaths:
      - stdout`

	confFile, _ := ioutil.TempFile(".", "config-*.yml")
	defer func() {
		confFile.Close()
		os.Remove(confFile.Name())
	}()

	confFile.WriteString(content)
	confFile.Sync()

	old := log.GetLogger()
	defer func() {
		log.SetLogger(old)
	}()

	l, err := log.NewLogger(confFile.Name())
	So(err, ShouldBeNil)
	log.SetLogger(l)

	app := fx.New(
		fx.Logger(log.GetLogger().CreateStdLogger()),
		marsconfig.NewModule(confFile.Name()),
	)

	So(test.StartFxApp(app), ShouldBeNil)
	defer test.StopFxApp(app)
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

type fakeRestServer struct {
	server      *echo.Echo
	Middlewares map[string]IMiddleware
}

func (r *fakeRestServer) BindMiddleware(m IMiddleware) {
	r.Middlewares[m.GetName()] = m
	r.server.Use(m.Handle)
}

func (r *fakeRestServer) SetRateLimitConfig(maxRequests, requestsPerSec int64) error {

	params := map[string]interface{}{}
	params["name"] = RateLimitName
	params["maxrequests"] = maxRequests
	params["requestspersec"] = requestsPerSec

	if middle, ok := r.Middlewares[RateLimitName]; ok {
		return middle.SetConfig(params)
	}

	return errors.New("not find ratelimt Middleware")
}

// Pre adds middleware to the chain which is run before router.
func (r *fakeRestServer) Pre(middleware ...echo.MiddlewareFunc) {
	r.server.Pre(middleware...)
}

// Use adds middleware to the chain which is run after router.
func (r *fakeRestServer) Use(middleware ...echo.MiddlewareFunc) {
	r.server.Use(middleware...)
}

func newFakeRestServer(server *echo.Echo) *fakeRestServer {
	return &fakeRestServer{
		server:      server,
		Middlewares: make(map[string]IMiddleware),
	}
}

func TestMiddlewareNotConfigured(t *testing.T) {

	Convey("Given echo and a other middleware config", t, func() {
		setLogger()
		confContent := `logger:
    name: logger
    order: 1`

		mwConf := make(map[string]interface{})
		yml.Unmarshal([]byte(confContent), mwConf)

		e := echo.New()
		e.POST("/hello", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		server := newFakeRestServer(e)

		boots := map[string]IMiddlewareBootstraper{
			"bodylimit":         &BodyLimitMiddlewareBootstraper{},
			"cors":              &CORSMiddlewareBootstraper{},
			"csrf":              &CSRFMiddlewareBootstraper{},
			"header":            &HeaderMiddlewareBootstraper{},
			"nonGetInterceptor": &NonGetInterceptorBootstraper{Availiable: true},
			"recover":           &RecoverMiddlewareBootstraper{},
			"requestid":         &RequestIDMiddlewareBootstraper{},
			"respcache":         &RespCacheMiddlewareBootstraper{},
			"stats":             &StatsMiddlewareBootstraper{},
		}

		for key, boot := range boots {
			boot.BindMiddleware(server, mwConf)

			whenCondition := fmt.Sprintf("When post a request without middleware %s", key)

			Convey(whenCondition, func() {
				hw := []byte("Hello, World!")
				status, resp := sendRequest(http.MethodPost, "/hello", hw, e)

				Convey("Then status should be http.StatusOK and response should be 'Hello, World!'", func() {
					So(status, ShouldEqual, http.StatusOK)
					So(string(resp), ShouldEqual, "Hello, World!")
				})
			})
		}
	})
}
