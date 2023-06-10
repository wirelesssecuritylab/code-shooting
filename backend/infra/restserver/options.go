package restserver

import (
	"net/http"
)

type Option interface {
	apply(*RestServer) error
}

type optionImpl struct {
	name    string
	optFunc func(*http.Server) error
}

func (s *optionImpl) apply(rs *RestServer) error {
	if rs.Name != s.name {
		return nil
	}

	return s.optFunc(rs.httpServer)
}

func WithOption(name string, optFunc func(*http.Server) error) Option {
	return &optionImpl{name, optFunc}
}
