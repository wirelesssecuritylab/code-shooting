package restserver

import (
	stdContext "context"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/net/netutil"

	middleconfig "code-shooting/infra/restserver/config"
	"code-shooting/infra/restserver/internal"
	"code-shooting/infra/restserver/middleware"
)

type (
	Context = echo.Context

	Router = echo.Router

	Route = echo.Route

	MiddlewareFunc = echo.MiddlewareFunc

	HandlerFunc = echo.HandlerFunc

	Group = echo.Group
)

type RestServer struct {
	sync.RWMutex
	Name         string
	server       *echo.Echo
	httpServer   *http.Server
	listener     net.Listener
	RootGroupBox *RootGroupBox
	Middlewares  map[string]middleware.IMiddleware
}

type RootGroupBox struct {
	RootGroup *echo.Group
	RootPath  string
}

// New creates an instance of RestServer.
func newRestServer(server *echo.Echo, conf *internal.RestServerConf, opts ...Option) (*RestServer, error) {
	httpserver, err := buildHttpServer(&conf.HttpServer)
	if err != nil {
		return nil, err
	}
	listener, err := buildListener(&conf.Listener)
	if err != nil {
		return nil, err
	}
	rootGroupBox := &RootGroupBox{}
	if conf.RootPath != "" {
		group := server.Group(conf.RootPath)
		rootGroupBox.RootGroup = group
		rootGroupBox.RootPath = conf.RootPath
	}
	rs := &RestServer{
		Name:         conf.Name,
		server:       server,
		httpServer:   httpserver,
		listener:     listener,
		RootGroupBox: rootGroupBox,
		Middlewares:  make(map[string]middleware.IMiddleware),
	}
	rs.useOptions(opts...)
	return rs, nil
}

func buildHttpServer(conf *internal.HttpServerConf) (*http.Server, error) {
	server := &http.Server{
		Addr:              conf.Addr,
		ReadTimeout:       conf.ReadTimeout,
		ReadHeaderTimeout: conf.ReadHeaderTimeout,
		WriteTimeout:      conf.WriteTimeout,
		IdleTimeout:       conf.IdleTimeout,
		MaxHeaderBytes:    conf.MaxHeaderBytes,
	}

	server.TLSConfig = nil
	if conf.Protocol == "https" && conf.CertFile != "" && conf.KeyFile != "" {
		cert, err := ioutil.ReadFile(conf.CertFile)
		if err != nil {
			return nil, errors.Wrap(err, "read TLS cert file")
		}
		key, err := ioutil.ReadFile(conf.KeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "read TLS key file")
		}
		certs, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, errors.Wrap(err, "TLS X509KeyPair")
		}
		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{certs}}
	}
	return server, nil
}

func buildListener(conf *internal.ListenerConf) (net.Listener, error) {

	listener, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		return nil, errors.Wrap(err, "listen tcp addr")
	}

	if conf.MaxConnections > 0 {
		listener = netutil.LimitListener(listener, conf.MaxConnections)
	}
	return listener, nil
}

func (r *RestServer) useOptions(opts ...Option) {
	for _, opt := range opts {
		opt.apply(r)
	}
}

// NewContext returns a Context instance.
func (r *RestServer) NewContext(req *http.Request, resp http.ResponseWriter) Context {
	return r.server.NewContext(req, resp)
}

// Router returns the default router.
func (r *RestServer) Router() *echo.Router {
	return r.server.Router()
}

func (r *RestServer) BindMiddleware(m middleware.IMiddleware) {
	if len(m.GetName()) > 0 {
		r.Lock()
		defer r.Unlock()
		r.Middlewares[m.GetName()] = m
	}

	r.server.Use(m.Handle)
}

func (r *RestServer) SetRateLimitConfig(maxRequests, requestsPerSec int64) error {
	params := map[string]interface{}{}
	params["name"] = middleware.RateLimitName
	params["maxrequests"] = maxRequests
	params["requestspersec"] = requestsPerSec

	r.Lock()
	defer r.Unlock()

	if middle, ok := r.Middlewares[middleware.RateLimitName]; ok {
		return middle.SetConfig(params)
	}

	return errors.New("not find ratelimt Middleware")
}

func (r *RestServer) updateMiddlewares(middlewaresConfig middleconfig.MiddlewareConf) error {
	errMsg := ""
	r.Lock()
	defer r.Unlock()

	for name, config := range middlewaresConfig {
		if middle, ok := r.Middlewares[name]; ok {
			if err := middle.SetConfig(config); err != nil {
				errMsg = errMsg + name + " has error: " + err.Error() + "\n"
			}
		}
	}
	if len(errMsg) == 0 {
		return nil
	}
	return errors.New(errMsg)
}

// Routers returns the map of host => router.
// func (r *RestServer) Routers() map[string]Router {
// 	return r.server.Routers()
// }

// DefaultHTTPErrorHandler is the default HTTP error handler. It sends a JSON response
// with status code.
func (r *RestServer) DefaultHTTPErrorHandler(err error, c Context) {
	r.server.DefaultHTTPErrorHandler(err, c)
}

// Pre adds middleware to the chain which is run before router.
func (r *RestServer) Pre(middleware ...MiddlewareFunc) {
	r.server.Pre(middleware...)
}

// Use adds middleware to the chain which is run after router.
func (r *RestServer) Use(middleware ...MiddlewareFunc) {
	r.server.Use(middleware...)
}

// CONNECT registers a new CONNECT route for a path with matching handler in the
// router with optional route-level middleware.
func (r *RestServer) CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return r.server.CONNECT(path, h, m...)
}

// DELETE registers a new DELETE route for a path with matching handler in the router
// with optional route-level middleware.
func (r *RestServer) DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return r.server.DELETE(path, h, m...)
}

// GET registers a new GET route for a path with matching handler in the router
// with optional route-level middleware.
func (r *RestServer) GET(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return r.server.GET(path, h, m...)
}

// HEAD registers a new HEAD route for a path with matching handler in the
// router with optional route-level middleware.
func (r *RestServer) HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return r.server.HEAD(path, h, m...)
}

// OPTIONS registers a new OPTIONS route for a path with matching handler in the
// router with optional route-level middleware.
func (r *RestServer) OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return r.server.OPTIONS(path, h, m...)
}

// PATCH registers a new PATCH route for a path with matching handler in the
// router with optional route-level middleware.
func (r *RestServer) PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return r.server.PATCH(path, h, m...)
}

// POST registers a new POST route for a path with matching handler in the
// router with optional route-level middleware.
func (r *RestServer) POST(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return r.server.POST(path, h, m...)
}

// PUT registers a new PUT route for a path with matching handler in the
// router with optional route-level middleware.
func (r *RestServer) PUT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return r.server.PUT(path, h, m...)
}

// TRACE registers a new TRACE route for a path with matching handler in the
// router with optional route-level middleware.
func (r *RestServer) TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return r.server.TRACE(path, h, m...)
}

// Any registers a new route for all HTTP methods and path with matching handler
// in the router with optional route-level middleware.
func (r *RestServer) Any(path string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route {
	return r.server.Any(path, handler, middleware...)
}

// Match registers a new route for multiple HTTP methods and path with matching
// handler in the router with optional route-level middleware.
func (r *RestServer) Match(methods []string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route {
	return r.server.Match(methods, path, handler, middleware...)
}

// Static registers a new route with path prefix to serve static files from the
// provided root directory.
func (r *RestServer) Static(prefix, root string) *Route {
	return r.server.Static(prefix, root)
}

// File registers a new route with path to serve a static file with optional route-level middleware.
func (r *RestServer) File(path, file string, m ...MiddlewareFunc) *Route {
	return r.server.File(path, file, m...)
}

// Add registers a new route for an HTTP method and path with matching handler
// in the router with optional route-level middleware.
func (r *RestServer) Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
	return r.server.Add(method, path, handler, middleware...)
}

// Host creates a new router group for the provided host and optional host-level middleware.
// func (r *RestServer) Host(name string, m ...MiddlewareFunc) (g *Group) {
// 	return r.server.Host(name, m...)
// }

// Group creates a new router group with prefix and optional group-level middleware.
func (r *RestServer) Group(prefix string, m ...MiddlewareFunc) (g *Group) {
	return r.server.Group(prefix, m...)
}

// URI generates a URI from handler.
func (r *RestServer) URI(handler HandlerFunc, params ...interface{}) string {
	return r.server.URI(handler, params...)
}

// URL is an alias for `URI` function.
func (r *RestServer) URL(h HandlerFunc, params ...interface{}) string {
	return r.server.URL(h, params...)
}

// Reverse generates an URL from route name and provided parameters.
func (r *RestServer) Reverse(name string, params ...interface{}) string {
	return r.server.Reverse(name, params...)
}

// Routes returns the registered routes.
func (r *RestServer) Routes() []*Route {
	return r.server.Routes()
}

// AcquireContext returns an empty `Context` instance from the pool.
// You must return the context by calling `ReleaseContext()`.
func (r *RestServer) AcquireContext() Context {
	return r.server.AcquireContext()
}

// ReleaseContext returns the `Context` instance back to the pool.
// You must call it after `AcquireContext()`.
func (r *RestServer) ReleaseContext(c Context) {
	r.server.ReleaseContext(c)
}

// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
func (r *RestServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	r.server.ServeHTTP(resp, req)
}

// Start starts an HTTP server.
func (r *RestServer) Start(address string) error {
	return r.server.Start(address)
}

// StartTLS starts an HTTPS server.
// If `certFile` or `keyFile` is `string` the values are treated as file paths.
// If `certFile` or `keyFile` is `[]byte` the values are treated as the certificate or key as-is.
func (r *RestServer) StartTLS(address string, certFile string, keyFile string) (err error) {
	return r.server.StartTLS(address, certFile, keyFile)
}

func (r *RestServer) StartAutoTLS(address string) error {
	return r.server.StartAutoTLS(address)
}

// StartServer starts a custom http server.
func (r *RestServer) StartServer(s *http.Server) (err error) {
	return r.server.StartServer(s)
}

// func (r *RestServer) StartH2CServer(address string, h2s *http2.server) (err error) {
// 	return r.server.StartH2CServer(address, h2s)
// }

func (r *RestServer) Close() error {
	return r.server.Close()
}

func (r *RestServer) Shutdown(ctx stdContext.Context) error {
	return r.server.Shutdown(ctx)
}
