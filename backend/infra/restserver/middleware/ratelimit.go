package middleware

import (
	"math"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"

	"code-shooting/infra/logger"
	middlewareconfig "code-shooting/infra/restserver/config"
)

var RateLimitName string = "ratelimit"

type RateLimitMiddlewareBootstraper struct {
	RateLimitMWConfig *RateLimitMiddlewareConfig
}

type RateLimitMiddlewareConfig struct {
	Name            string `yaml:"name"`
	Order           int    `yaml:"order"`
	RateLimitConfig `yaml:",inline"`
}

type RateLimitConfig struct {
	RequestsPerSec int64 `yaml:"requestspersec"`
	MaxRequests    int64 `yaml:"maxrequests"`
	limiter        *rate.Limiter
}

var _ IMiddlewareBootstraper = &RateLimitMiddlewareBootstraper{}

func (boot *RateLimitMiddlewareBootstraper) Enable(config middlewareconfig.MiddlewareConf) bool {

	_, ok := config["ratelimit"]
	return ok
}

func (boot *RateLimitMiddlewareBootstraper) Order(config middlewareconfig.MiddlewareConf) int {

	var order = COMMON_ORDER

	if ratelimit, ok := config["ratelimit"]; ok {
		boot.RateLimitMWConfig = &RateLimitMiddlewareConfig{}
		if err := transformInterfaceToObject(ratelimit, boot.RateLimitMWConfig); err == nil {
			order = boot.RateLimitMWConfig.Order + order
		}
	}
	return order
}

func (boot *RateLimitMiddlewareBootstraper) BindMiddleware(server RestServer, config middlewareconfig.MiddlewareConf) {

	boot.RateLimitMWConfig = &RateLimitMiddlewareConfig{}

	ratelimit, ok := config["ratelimit"]
	if !ok {
		logger.Info("the RateLimit Middleware is not configured, ignore")
		return
	}
	logger.Infof("enable the RateLimit Middleware, config is %v", ratelimit)

	if err := transformInterfaceToObject(ratelimit, boot.RateLimitMWConfig); err != nil {
		logger.Warnf("the RateLimit Middleware is error config: %s, ignore", err.Error())
		return
	}

	var m IMiddleware = &RateLimitMiddleware{}
	m.SetConfig(ratelimit)
	server.BindMiddleware(m)
}

func (boot *RateLimitMiddlewareBootstraper) RateLimitWithConfig(config RateLimitConfig) echo.MiddlewareFunc {

	if config.limiter == nil {
		config.limiter = rate.NewLimiter(rate.Limit(config.RequestsPerSec), int(config.MaxRequests))
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !config.limiter.Allow() {
				c.Response().Header().Add("Retry-After", boot.getRetryAfter(config))
				return echo.ErrTooManyRequests
			}
			return next(c)
		}
	}
}

func (boot *RateLimitMiddlewareBootstraper) getRetryAfter(config RateLimitConfig) string {
	return strconv.Itoa(int(math.Ceil(float64(config.MaxRequests / config.RequestsPerSec))))
}

type RateLimitMiddleware struct {
	sync.RWMutex
	RequestsPerSec int64
	MaxRequests    int64
	limiter        *rate.Limiter
}

func (s *RateLimitMiddleware) GetName() string {
	return RateLimitName
}

func (s *RateLimitMiddleware) SetConfig(params interface{}) error {
	cfg := &RateLimitMiddlewareConfig{}
	if err := transformInterfaceToObject(params, cfg); err != nil {
		logger.Warnf("the RateLimit Middleware is error config: %s, ignore", err.Error())
		return err
	}

	if s.isParamsValid(cfg.MaxRequests, cfg.RequestsPerSec) {

		if cfg.RequestsPerSec > cfg.MaxRequests {
			logger.Warnf("the RateLimitMiddleware : tokens per sec（%d） is more than max tokens(%d) !!!",
				cfg.RequestsPerSec, cfg.MaxRequests)
		}

		s.Lock()
		defer s.Unlock()
		s.MaxRequests = cfg.MaxRequests
		s.RequestsPerSec = cfg.RequestsPerSec
		if s.limiter == nil {
			s.limiter = rate.NewLimiter(rate.Limit(s.RequestsPerSec), int(s.MaxRequests))
		} else {
			s.limiter.SetBurst(int(cfg.MaxRequests))
			s.limiter.SetLimit(rate.Limit(cfg.RequestsPerSec))
		}
	} else {
		logger.Warnf("the RateLimit Middleware Param is error: MaxRequests %d, RequestsPerSec %d, ignore",
			cfg.MaxRequests, cfg.RequestsPerSec)
	}

	return nil
}

func (s *RateLimitMiddleware) isParamsValid(maxRequests, requestsPerSec int64) bool {
	s.RLock()
	defer s.RUnlock()
	if maxRequests > 0 && requestsPerSec > 0 && !(s.MaxRequests == maxRequests && s.RequestsPerSec == requestsPerSec) {
		return true
	}
	return false
}

func (s *RateLimitMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s.RLock()
		defer s.RUnlock()

		if !s.limiter.Allow() {
			c.Response().Header().Add("Retry-After", s.getRetryAfter())
			return echo.ErrTooManyRequests
		}
		return next(c)
	}
}

func (s *RateLimitMiddleware) getRetryAfter() string {
	return strconv.Itoa(int(math.Ceil(float64(s.MaxRequests / s.RequestsPerSec))))
}
