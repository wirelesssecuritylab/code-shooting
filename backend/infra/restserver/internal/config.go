package internal

import (
	"code-shooting/infra/restserver/config"
	"time"
)

type MultRestServerConf = map[string]*RestServerConf

type RestServerConf struct {
	Name        string
	HttpServer  HttpServerConf
	Listener    ListenerConf
	RootPath    string
	Middlewares config.MiddlewareConf
}
type HttpServerConf struct {
	Protocol          string
	Addr              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
	CertFile          string
	KeyFile           string
}

type ListenerConf struct {
	Addr           string
	MaxConnections int
}
