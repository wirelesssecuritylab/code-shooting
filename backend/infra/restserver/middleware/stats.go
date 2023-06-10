package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/thoas/stats"

	"code-shooting/infra/logger"
	middlewareconfig "code-shooting/infra/restserver/config"
)

type StatsMiddlewareBootstraper struct {
	StatsConfig *StatsMiddlewareConfig
	Middleware  *stats.Stats
	RecordChan  chan *statsRecord
	ClearTime   time.Time
}

type StatsMiddlewareConfig struct {
	Name  string `yaml:"name"`
	Order int    `yaml:"order"`
}

type statsRecord struct {
	statusCode   int
	start        time.Time
	responseSize int64
}

var _ IMiddlewareBootstraper = &StatsMiddlewareBootstraper{}

func (boot *StatsMiddlewareBootstraper) Enable(config middlewareconfig.MiddlewareConf) bool {

	_, ok := config["stats"]
	return ok
}

func (boot *StatsMiddlewareBootstraper) Order(config middlewareconfig.MiddlewareConf) int {

	var order = COMMON_ORDER
	if stats, ok := config["stats"]; ok {

		boot.StatsConfig = &StatsMiddlewareConfig{}

		if err := transformInterfaceToObject(stats, boot.StatsConfig); err == nil {
			order = boot.StatsConfig.Order
		}
	}
	return order + COMMON_ORDER
}

func (boot *StatsMiddlewareBootstraper) BindMiddleware(server RestServer, config middlewareconfig.MiddlewareConf) {

	boot.StatsConfig = &StatsMiddlewareConfig{}
	boot.RecordChan = make(chan *statsRecord, 1024)
	boot.ClearTime = time.Now()

	stats, ok := config["stats"]
	if !ok {
		//用户没有配置
		logger.Info("the Stats Middleware is not configured, ignore")
		return
	}
	logger.Infof("enable the Stats Middleware, config is %v", stats)

	server.Use(boot.createstats())
	go boot.end()

}

func (boot *StatsMiddlewareBootstraper) createstats() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			beginning, _ := boot.Middleware.Begin(c.Response())

			err := next(c)

			save := !isWebsocketRequest(c.Request())
			if save {
				defer func() {
					record := &statsRecord{start: beginning, statusCode: c.Response().Status, responseSize: c.Response().Size}
					boot.RecordChan <- record
				}()
			}
			return err
		}
	}

}

func (boot *StatsMiddlewareBootstraper) end() {
	for record := range boot.RecordChan {
		//TODO 可能会有并发冲突的问题，在读取的时候，但是考虑只是一个运维功能，暂时不存在这个问题
		if time.Now().Sub(boot.ClearTime).Minutes() > 5 {
			boot.Middleware.ResponseCounts = make(map[string]int)
			boot.Middleware.TotalResponseCounts = make(map[string]int)
			boot.Middleware.TotalResponseTime = time.Time{}
			boot.Middleware.TotalResponseSize = int64(0)
			boot.ClearTime = time.Now()
		}

		responseTime := time.Since(record.start)
		statusCode := strconv.Itoa(record.statusCode)
		fmt.Println("response", boot.Middleware.ResponseCounts == nil)
		boot.Middleware.ResponseCounts[statusCode]++
		boot.Middleware.TotalResponseCounts[statusCode]++
		boot.Middleware.TotalResponseTime = boot.Middleware.TotalResponseTime.Add(responseTime)
		boot.Middleware.TotalResponseSize += record.responseSize

	}

}

func newOptions(options ...stats.Option) *stats.Options {
	opts := &stats.Options{}
	for _, o := range options {
		o(opts)
	}
	return opts
}

func isWebsocketRequest(req *http.Request) bool {
	return containsHeader(req, "Connection", "upgrade") && containsHeader(req, "Upgrade", "websocket")
}

func containsHeader(req *http.Request, name, value string) bool {
	items := strings.Split(req.Header.Get(name), ",")
	for _, item := range items {
		if value == strings.ToLower(strings.TrimSpace(item)) {
			return true
		}
	}
	return false
}
