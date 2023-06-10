package internal

import (
	"code-shooting/infra/restserver/config"
	"time"
)

type MultRestServerDTO []*RestServerDTO

type RestServerDTO struct {
	Name              string
	Protocol          string
	Addr              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
	MaxConnections    int
	CertFile          string
	KeyFile           string
	RootPath          string
	Middlewares       config.MiddlewareDTOs
}
