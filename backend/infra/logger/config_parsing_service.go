package logger

import (
	"regexp"
	"strconv"
	"sync"

	"code-shooting/infra/common"
	"code-shooting/infra/config"
	"code-shooting/infra/logger/internal"

	"github.com/pkg/errors"
)

type configParsingService interface {
	Parse(config config.Config, options ...Option) (internal.Config, error)
}

var _configParsingService configParsingService
var _configParsingServiceOnce sync.Once

func getConfigParsingService() configParsingService {
	_configParsingServiceOnce.Do(func() {
		_configParsingService = &configParsingServiceImpl{}
	})
	return _configParsingService
}

type configParsingServiceImpl struct {
}

func (s *configParsingServiceImpl) Parse(config config.Config, options ...Option) (internal.Config, error) {
	dto, err := s.buildDto(config)
	if err != nil {
		return internal.Config{}, errors.Wrap(err, "config format")
	}
	Options(options).Do(&dto)
	return s.transformToModel(dto)
}

func (s *configParsingServiceImpl) buildDto(config config.Config) (configDto, error) {
	if config == nil {
		return configDto{}, errors.New("go mars config is nil")
	}
	dto := &configDto{}
	err := config.Get(common.ROOT+"."+common.LOG, dto)
	return *dto, err
}

func (s *configParsingServiceImpl) transformToModel(dto configDto) (internal.Config, error) {
	cfg := internal.Config{}

	cfg.Level = internal.Level(dto.Level)
	err := cfg.Level.Check()
	if err != nil {
		return internal.Config{}, errors.Wrap(err, "check level")
	}

	cfg.Encoder = internal.Encoder(dto.Encoder)
	err = cfg.Encoder.Check()
	if err != nil {
		return internal.Config{}, errors.Wrap(err, "check encoder")
	}

	format, err := s.transformCustomFormat(dto.Format, cfg.Encoder)
	if err != nil {
		return internal.Config{}, errors.Wrap(err, "check custom format")
	}
	cfg.Format = format

	cfg.OutputPaths = dto.OutputPaths
	if len(cfg.OutputPaths) == 0 {
		return internal.Config{}, errors.New("output paths is empty")
	}

	cfg.RotateConfig, err = s.transformRotateConfig(dto.RotateConfig)
	if err != nil {
		return internal.Config{}, errors.Wrap(err, "check rotateConfig")
	}

	return cfg, nil
}

func (s *configParsingServiceImpl) transformCustomFormat(format string, encoder internal.Encoder) (internal.CustomFormat, error) {
	if encoder.Equal(internal.PlainEncoder) {
		format, err := internal.GetCustomFormatParsingService().Parse(format)
		if err != nil {
			return internal.CustomFormat{}, errors.Wrap(err, "check custom format of plain encoder")
		}
		return format, nil
	}
	return internal.CustomFormat{}, nil
}

func (s *configParsingServiceImpl) normalizeRotateConfig(config internal.RotateConfig) (internal.RotateConfig, error) {
	cfg := config
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = internal.DefaultMaxSize
	}

	if cfg.MaxBackups <= 0 && cfg.MaxAge <= 0 {
		cfg.MaxBackups = internal.DefaultMaxBackups
	}

	if cfg.FileMode == "" {
		cfg.FileMode = internal.DefaultFileMode
	}

	match, _ := regexp.MatchString("^0[0-7]{3}$", cfg.FileMode)
	if match == false {
		return internal.RotateConfig{}, errors.New("fileMode is not octal number system")
	}
	_, err := strconv.ParseInt(cfg.FileMode, 8, 0)
	if err != nil {
		return internal.RotateConfig{}, errors.Wrap(err, "check fileMode for strconv.ParseInt")
	}

	return cfg, nil
}

func (s *configParsingServiceImpl) transformRotateConfig(rotate *rotateConfigDto) (internal.RotateConfig, error) {
	cfg := internal.RotateConfig{}

	if rotate == nil {
		cfg.Disable = true
	} else {
		cfg.MaxSize = rotate.MaxSize
		cfg.MaxBackups = rotate.MaxBackups
		cfg.MaxAge = rotate.MaxAge
		cfg.Compress = rotate.Compress
		cfg.FileMode = rotate.FileMode
	}

	return s.normalizeRotateConfig(cfg)
}
