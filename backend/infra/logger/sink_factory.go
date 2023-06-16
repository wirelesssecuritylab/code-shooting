package logger

import (
	"code-shooting/infra/lumberjack"
	"net/url"
	"strconv"
	"sync"

	"code-shooting/infra/logger/internal"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type sinkFactory interface {
	createSink(outputPath string, rotateCfg internal.RotateConfig) (*writeSyncer, error)
}

var _sinkFactory sinkFactory
var _sinkFactoryOnce sync.Once

func getSinkFactory() sinkFactory {
	_sinkFactoryOnce.Do(func() {
		_sinkFactory = &sinkFactoryImpl{}
	})
	return _sinkFactory
}

type writeSyncer struct {
	syncer zapcore.WriteSyncer
	closer func()
}

type sinkFactoryImpl struct {
}

func (s *sinkFactoryImpl) createSink(outputPath string, rotateCfg internal.RotateConfig) (*writeSyncer, error) {
	if !s.needRotate(outputPath, rotateCfg) {
		sink, closer, err := zap.Open(outputPath)
		if err != nil {
			return nil, errors.Wrapf(err, "create sink of %s", outputPath)
		}
		return &writeSyncer{sink, closer}, err
	}

	mode, _ := strconv.ParseInt(rotateCfg.FileMode, 8, 0)
	lLogger := &lumberjack.Logger{
		Filename:   outputPath,
		MaxSize:    rotateCfg.MaxSize,
		MaxAge:     rotateCfg.MaxAge,
		MaxBackups: rotateCfg.MaxBackups,
		Compress:   rotateCfg.Compress,
		FileMode:   uint32(mode),
	}
	close := func() {
		err := lLogger.Close()
		if err != nil {
			Error("close lumberjack logger err : ", err.Error())
		}
	}
	return &writeSyncer{zapcore.AddSync(lLogger), close}, nil
}

func (s *sinkFactoryImpl) needRotate(outputPath string, config internal.RotateConfig) bool {
	if config.Disable {
		return false
	}

	noNeedFiles := internal.NewStringSet()
	noNeedFiles.Add("stdout", "stderr")

	return !s.isRegistedScheme(outputPath) && !noNeedFiles.Contains(outputPath)
}

func (s *sinkFactoryImpl) isRegistedScheme(outputPath string) bool {
	u, err := url.Parse(outputPath)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Scheme != "file"
}
