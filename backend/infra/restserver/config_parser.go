package restserver

import (
	"code-shooting/infra/common/constants"
	"net"

	"code-shooting/infra/restserver/internal"

	marsconfig "code-shooting/infra/config"
	middlewareconfig "code-shooting/infra/restserver/config"
	"github.com/pkg/errors"
)

type configParsingService interface {
	Parse(config marsconfig.Config) (internal.MultRestServerConf, error)
}

func GetConfigParser() configParsingService {
	return &configParser{}
}

type configParser struct {
}

func (s *configParser) Parse(config marsconfig.Config) (internal.MultRestServerConf, error) {
	multDTO, err := s.buildMultDTO(config)
	if err != nil {
		return nil, errors.Wrap(err, "config format")
	}
	return s.transformToModels(multDTO)
}

func (s *configParser) buildMultDTO(config marsconfig.Config) (internal.MultRestServerDTO, error) {
	if config == nil {
		return nil, errors.New("config does not exist")
	}

	var multDTO internal.MultRestServerDTO
	err := config.Get(constants.ROOT+"."+constants.REST_SERVER, &multDTO)
	if err != nil {
		return nil, errors.Wrap(err, "get restservers config")
	}
	return multDTO, nil
}

func (s *configParser) transformToModels(multDTO internal.MultRestServerDTO) (internal.MultRestServerConf, error) {

	multModel := make(map[string]*internal.RestServerConf)
	for _, dto := range multDTO {
		if err := s.check(dto); err != nil {
			return nil, errors.Wrap(err, "check dto")
		}

		model := s.transformToModel(dto)
		_, ok := multModel[model.Name]
		if ok {
			return nil, errors.Errorf("name %s is duplicated", model.Name)
		}

		multModel[model.Name] = model
	}
	return multModel, nil
}

func (s *configParser) check(dto *internal.RestServerDTO) error {
	host, _, err := net.SplitHostPort(dto.Addr)
	if err != nil {
		return errors.Wrap(err, "split host port from address")
	}
	if ip := net.ParseIP(host); ip == nil {
		return errors.New("not a valid textual representation of an IP address")
	}
	return nil
}

func (s *configParser) transformToModel(dto *internal.RestServerDTO) *internal.RestServerConf {

	return &internal.RestServerConf{
		Name:     dto.Name,
		RootPath: dto.RootPath,
		HttpServer: internal.HttpServerConf{
			Protocol:          dto.Protocol,
			Addr:              dto.Addr,
			ReadTimeout:       dto.ReadTimeout,
			ReadHeaderTimeout: dto.ReadHeaderTimeout,
			WriteTimeout:      dto.WriteTimeout,
			IdleTimeout:       dto.IdleTimeout,
			MaxHeaderBytes:    dto.MaxHeaderBytes,
			CertFile:          dto.CertFile,
			KeyFile:           dto.KeyFile,
		},
		Listener: internal.ListenerConf{
			Addr:           dto.Addr,
			MaxConnections: dto.MaxConnections,
		},
		Middlewares: s.transformToMiddlewareMap(dto.Middlewares),
	}
}

func (s *configParser) transformToMiddlewareMap(middlewares middlewareconfig.MiddlewareDTOs) middlewareconfig.MiddlewareConf {

	middlewareMap := make(map[string]interface{})
	order := 0
	for _, mw := range middlewares {
		if name, ok := mw["name"].(string); ok {
			order = order + 1
			mw["order"] = order
			middlewareMap[name] = mw
		}
	}
	return middlewareMap

}
