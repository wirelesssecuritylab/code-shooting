package sql

import (
	"database/sql"
	"time"

	"code-shooting/infra/common"
	marsconfig "code-shooting/infra/config"

	"github.com/pkg/errors"
)

const (
	_MaxOpenConns    = "maxOpen"
	_MaxIdleConns    = "maxIdle"
	_ConnMaxLifetime = "maxLifetime"
)

type configParsingService interface {
	Parse(config marsconfig.Config) (map[string]config, error)
}

func getConfigParsingService() configParsingService {
	return &configParsingServiceImpl{}
}

type configParsingServiceImpl struct {
}

func (s *configParsingServiceImpl) Parse(config marsconfig.Config) (map[string]config, error) {
	dtos, err := s.buildDTOs(config)
	if err != nil {
		return nil, errors.Wrap(err, "config format")
	}
	return s.transformToModels(dtos)
}

func (s *configParsingServiceImpl) buildDTOs(config marsconfig.Config) ([]configDTO, error) {
	if config == nil {
		return nil, errors.New("config does not exist")
	}

	var dtos []configDTO
	err := config.Get(common.ROOT+"."+common.DATABASE_PG, &dtos)
	if err != nil {
		return nil, errors.Wrap(err, "query pg config")
	}
	return dtos, nil
}

func (s *configParsingServiceImpl) transformToModels(dtos []configDTO) (map[string]config, error) {
	configs := make(map[string]config)
	//configs := make([]config, 0, len(dtos))
	for _, dto := range dtos {
		cfg, err := s.transformToModel(dto)
		if err != nil {
			return nil, errors.Wrapf(err, "transform to model of %v", dto)
		}
		_, ok := configs[cfg.ID]
		if ok {
			return nil, errors.Errorf("id %s is duplicated", cfg.ID)
		}
		configs[cfg.ID] = cfg
	}
	return configs, nil
}

func (s *configParsingServiceImpl) transformToModel(dto configDTO) (config, error) {
	cfg := config{
		ID:         dto.ID,
		User:       dto.User,
		Password:   dto.Password,
		Host:       dto.Host,
		Port:       dto.Port,
		DBName:     dto.DBName,
		ConnParams: make(map[string]string),
	}
	err := s.fillParams(&cfg, &dto)
	if err != nil {
		return config{}, errors.Wrap(err, "fill params")
	}
	err = s.fillPlugins(&cfg, &dto)
	return cfg, err
}

func (s *configParsingServiceImpl) fillPlugins(cfg *config, dto *configDTO) error {
	if len(dto.Password) > 0 {
		return nil
	}
	if len(dto.Plugins) != 1 {
		return errors.Errorf("just support one plugin(%d)", len(dto.Plugins))
	}

	return nil
}

func (s *configParsingServiceImpl) fillParams(cfg *config, dto *configDTO) error {
	filler := &dbOperationParamFiller{}
	filler.setNextFiller(&connParamFiller{})

	for k, v := range dto.ConnParams {
		err := filler.fill(cfg, k, v)
		if err != nil {
			return errors.Wrapf(err, "fill params of %s:%v", k, v)
		}
	}
	return nil
}

type paramFiller interface {
	fill(cfg *config, key string, value interface{}) error
	setNextFiller(filler paramFiller)
}

type dbOperationParamFiller struct {
	next paramFiller
}

func (s *dbOperationParamFiller) fill(cfg *config, key string, value interface{}) error {
	v, ok := value.(int)

	wrap := func(do func()) error {
		if !ok {
			return errors.Errorf("value(%v) type of %s is not int", value, key)
		}
		do()
		return nil
	}

	switch key {
	case _MaxOpenConns:
		return wrap(func() {
			cfg.DBOptions = append(cfg.DBOptions, func(db *sql.DB) {
				db.SetMaxOpenConns(v)
			})
		})
	case _MaxIdleConns:
		return wrap(func() {
			cfg.DBOptions = append(cfg.DBOptions, func(db *sql.DB) {
				db.SetMaxIdleConns(v)
			})
		})
	case _ConnMaxLifetime:
		return wrap(func() {
			cfg.DBOptions = append(cfg.DBOptions, func(db *sql.DB) {
				db.SetConnMaxLifetime(time.Duration(v) * time.Second)
			})
		})
	default:
		if s.next == nil {
			return errors.Errorf("can not find filler handle %s of db operation filler", key)
		}
		return s.next.fill(cfg, key, value)
	}
}

func (s *dbOperationParamFiller) setNextFiller(filler paramFiller) {
	s.next = filler
}

type connParamFiller struct {
	next paramFiller
}

func (s *connParamFiller) fill(cfg *config, key string, value interface{}) error {
	switch v := value.(type) {
	case string:
		cfg.ConnParams[key] = v
	default:
		if s.next == nil {
			return errors.Errorf("can not find filler handle %s of conn param filler", key)
		}
		return s.next.fill(cfg, key, value)
	}
	return nil
}

func (s *connParamFiller) setNextFiller(filler paramFiller) {
	s.next = filler
}
