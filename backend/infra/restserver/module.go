package restserver

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"code-shooting/infra/restserver/internal"

	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/thoas/stats"
	"go.uber.org/fx"

	"code-shooting/infra/logger"
	middlewareconfig "code-shooting/infra/restserver/config"
	"code-shooting/infra/restserver/expvar"
	"code-shooting/infra/restserver/middleware"
	"code-shooting/infra/restserver/pprof"
	"code-shooting/infra/restserver/utils/zapwrap"
	"code-shooting/infra/restserver/validator"
)

var banner = `
  ____             __  __                
 / ___| ___       |  \/  | __ _ _ __ ___ 
| |  _ / _ \ _____| |\/| |/ _| | '__/ __|
| |_| | (_) |_____| |  | | (_| | |  \__ \
 \____|\___/      |_|  |_|\__,_|_|  |___/      v1.5.1
_____________________________________________________

`

type MiddleWareOutResult struct {
	fx.Out
	Bootstraper middleware.IMiddlewareBootstraper `group:"echomiddleware"`
}

func NewCORSBootstraper() MiddleWareOutResult {

	return MiddleWareOutResult{
		Bootstraper: &middleware.CORSMiddlewareBootstraper{},
	}
}

func NewCSRFBootstraper() MiddleWareOutResult {

	return MiddleWareOutResult{
		Bootstraper: &middleware.CSRFMiddlewareBootstraper{},
	}
}

func NewRespCacheBootstraper() MiddleWareOutResult {

	return MiddleWareOutResult{
		Bootstraper: &middleware.RespCacheMiddlewareBootstraper{},
	}
}

// func NewLoggerBootstraper() MiddleWareOutResult {

// 	return MiddleWareOutResult{
// 		Bootstraper: &middleware.LoggerMiddlewareBootstraper{},
// 	}
// }

func NewRecoverBootstraper() MiddleWareOutResult {

	return MiddleWareOutResult{
		Bootstraper: &middleware.RecoverMiddlewareBootstraper{},
	}
}

var statsMiddleWare = stats.New()

func NewStatsBootstraper() MiddleWareOutResult {

	return MiddleWareOutResult{
		Bootstraper: &middleware.StatsMiddlewareBootstraper{Middleware: statsMiddleWare},
	}
}

func NewRequestIDBootstraper() MiddleWareOutResult {

	return MiddleWareOutResult{
		Bootstraper: &middleware.RequestIDMiddlewareBootstraper{},
	}
}

func NewRequestHeaderBootstraper() MiddleWareOutResult {

	return MiddleWareOutResult{
		Bootstraper: &middleware.HeaderMiddlewareBootstraper{},
	}
}

func NewRequestMethodOverrideBootstraper() MiddleWareOutResult {

	return MiddleWareOutResult{
		Bootstraper: &middleware.MethodOverrideMiddlewareBootstraper{},
	}
}

// func NewLoggerErrorBootstraper() MiddleWareOutResult {

// 	return MiddleWareOutResult{
// 		Bootstraper: &middleware.LoggerErrorMiddlewareBootstraper{},
// 	}
// }

// type GrParams struct {
// 	fx.In
// 	Subject *gr.ServiceStateSubject `optional:"true"`
// }

// func NewNonGetInterceptorBootstraper(grParams GrParams) MiddleWareOutResult {

// 	bootStraper := &middleware.NonGetInterceptorBootstraper{}

// 	if grParams.Subject == nil {

// 		logger().Info("gr module is not provided! ")
// 		// 设置过滤器状态为不可用
// 		bootStraper.Availiable = false
// 	} else {
// 		// 设置过滤器状态为可用
// 		bootStraper.Availiable = true
// 		grParams.Subject.RegisterObserver(bootStraper)
// 		bootStraper.IsMaster = grParams.Subject.IsMaster()
// 	}

// 	return MiddleWareOutResult{
// 		Bootstraper: bootStraper,
// 	}
// }

func NewBodyLimitBootstraper() MiddleWareOutResult {
	return MiddleWareOutResult{
		Bootstraper: &middleware.BodyLimitMiddlewareBootstraper{},
	}
}

func NewRateLimitBootstraper() MiddleWareOutResult {
	return MiddleWareOutResult{
		Bootstraper: &middleware.RateLimitMiddlewareBootstraper{},
	}
}

type RestServerParam struct {
	fx.In
	Conf         *internal.RestServerConf
	Bootstrapers []middleware.IMiddlewareBootstraper `group:"echomiddleware"`
	Opts         []Option                            `optional:"true"`
}

func newRestServerModule(lc fx.Lifecycle, params RestServerParam) (*RestServer, error) {

	echo := echo.New()
	//设置校验器
	echo.Validator = &validator.Validator{}
	echo.HideBanner = true
	echo.Logger = zapwrap.Wrap()

	restserver, err := newRestServer(echo, params.Conf, params.Opts...)
	if err != nil {
		return nil, errors.Wrap(err, "new rest server")
	}

	bindMiddlewares(restserver, params)

	//在绑定的同时进行校验
	echo.Binder = validator.NewValidatorBinder()

	if params.Conf.RootPath != "" {
		g := echo.Group(params.Conf.RootPath)
		restserver.RootGroupBox.RootGroup = g
		restserver.RootGroupBox.RootPath = params.Conf.RootPath
	}

	lc.Append(fx.Hook{

		OnStart: func(ctx context.Context) error {

			for {
				hasError := false
				err := startRestServer(restserver, params, &hasError)
				if !hasError {
					log.Printf("\n%s\n", banner)
					logger.Info("restserver is started")
					return nil
				}
				log.Println("restserver start error: ", err.Error()) //console log

				select {
				case <-ctx.Done():
					logger.Error("restserver start failed: ", ctx.Err())
					log.Println("restserver start failed: ", ctx.Err()) //console log
					return errors.Wrap(ctx.Err(), "restserver start failed"+err.Error())
				default:
					logger.Warn("restserver retry start...")
					time.Sleep(5 * time.Second)
				}
			}

		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping restserver HTTP server: ", params.Conf.HttpServer.Addr)
			return stopRestServer(ctx, restserver)
		},
	})

	return restserver, nil

}

func startRestServer(restserver *RestServer, params RestServerParam, hasError *bool) (res error) {

	logger.Info("Starting restserver HTTP server: ", params.Conf.HttpServer.Addr)
	for _, route := range restserver.server.Routes() {
		logger.Info("restserver route item info: ", route)
	}

	go func() {

		if restserver.httpServer.TLSConfig == nil {
			restserver.server.Listener = restserver.listener
			restserver.server.TLSListener = nil
		} else {
			restserver.server.Listener = nil
			restserver.server.TLSListener = tls.NewListener(restserver.listener, restserver.httpServer.TLSConfig)
		}
		if err := restserver.server.StartServer(restserver.httpServer); err != nil {
			if err == http.ErrServerClosed {
				logger.Info("restserver closed")
			} else {
				*hasError = true
				res = errors.Wrap(err, "restserver start")
				logger.Warn("restserver start failed: ", err.Error())
			}
			return
		}
	}()
	time.Sleep(1 * time.Second)
	return

}

func stopRestServer(ctx context.Context, restserver *RestServer) error {

	if restserver.server.Listener != nil {
		err := restserver.server.Listener.Close()
		logger.Debug("restserver listener close err:", err)
	}
	if restserver.server.TLSListener != nil {
		err := restserver.server.TLSListener.Close()
		logger.Debug("restserver TLSListener close err:", err)
	}

	err := restserver.server.Shutdown(ctx)
	logger.Debug("restserver shutdown err:", err)
	err = restserver.httpServer.Shutdown(ctx)
	logger.Debug("http server shutdown err:", err)
	return err
}

func orderBootstrapers(bootstrapers []middleware.IMiddlewareBootstraper, middlewareConf middlewareconfig.MiddlewareConf) []middleware.IMiddlewareBootstraper {

	bootstraperMap := make(map[int][]middleware.IMiddlewareBootstraper)
	orders := make([]int, 0)
	for _, bootstrap := range bootstrapers {

		order := bootstrap.Order(middlewareConf)
		slice, ok := bootstraperMap[order]
		if !ok {
			slice = make([]middleware.IMiddlewareBootstraper, 0)
			bootstraperMap[order] = slice
			orders = append(orders, order)
		}
		bootstraperMap[order] = append(bootstraperMap[order], bootstrap)
	}

	sort.Ints(orders)

	orderBootstrpas := make([]middleware.IMiddlewareBootstraper, 0)
	for _, key := range orders {

		if key == middleware.COMMON_ORDER {
			// filter the middlewares that is not configured.
			continue
		}

		values := bootstraperMap[key]
		orderBootstrpas = append(orderBootstrpas, values...)
	}
	return orderBootstrpas

}

func bindMiddlewares(restserver *RestServer, params RestServerParam) {

	logger.Info("restserver middlewaresMap: ", params.Conf.Middlewares)

	boots := orderBootstrapers(params.Bootstrapers, params.Conf.Middlewares)
	for _, boot := range boots {
		boot.BindMiddleware(restserver, params.Conf.Middlewares)
	}

}

func RegisterExpVar(multRestServer *MultRestServer, multRestServerConf internal.MultRestServerConf) {

	for rsname, rsconf := range multRestServerConf {
		if _, ok := rsconf.Middlewares["expvar"]; ok {
			restserver := multRestServer.RestServerMap[rsname]
			expvar.RegisterExpVarHandler(restserver.server)
		}
	}
}

func RegisterPProf(multRestServer *MultRestServer, multRestServerConf internal.MultRestServerConf) {

	for rsname, rsconf := range multRestServerConf {
		if _, ok := rsconf.Middlewares["pprof"]; ok {
			restserver := multRestServer.RestServerMap[rsname]
			pprof.RegisterPProfHandler("", restserver.server)
		}
	}
}

func statsHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		stats := statsMiddleWare.Data()
		b, _ := json.Marshal(stats)
		w.Write(b)
	})

}

func RegisterStats(multRestServer *MultRestServer, multRestServerConf internal.MultRestServerConf) {

	for rsname, rsconf := range multRestServerConf {
		if _, ok := rsconf.Middlewares["stats"]; ok {
			restserver := multRestServer.RestServerMap[rsname]
			restserver.GET("/stats", echo.WrapHandler(statsHandler()))
		}
	}
}

func NewModule() fx.Option {

	return fx.Options(
		fx.Provide(NewCORSBootstraper),
		fx.Provide(NewCSRFBootstraper),
		fx.Provide(NewRespCacheBootstraper),
		fx.Provide(NewRecoverBootstraper),
		fx.Provide(NewRequestIDBootstraper),
		fx.Provide(NewStatsBootstraper),
		// fx.Provide(NewNonGetInterceptorBootstraper),
		fx.Provide(NewBodyLimitBootstraper),
		fx.Provide(NewRateLimitBootstraper),
		fx.Provide(NewRequestHeaderBootstraper),
		fx.Provide(NewRequestMethodOverrideBootstraper),
		fx.Provide(NewMultRestServerConf),
		fx.Provide(NewMultRestServer),
		fx.Invoke(RegisterExpVar),
		fx.Invoke(RegisterPProf),
		fx.Invoke(RegisterStats),
	)
}
