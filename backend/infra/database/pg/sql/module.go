package sql

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/fx"

	marsconfig "code-shooting/infra/config"
	"code-shooting/infra/logger"
)

type optionalHttpClientFactory struct {
	fx.In
}

func NewModule() fx.Option {
	return fx.Provide(func(lc fx.Lifecycle, config marsconfig.Config, httpClientFactory optionalHttpClientFactory) (SqlService, error) {
		sqlService, err := NewSqlService(config)
		if err != nil {
			return nil, err
		}
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				for {
					logger.Info("check pg service")
					err := sqlService.CheckStatus()
					if err == nil {
						logger.Info("pg service is ready")
						return nil
					}
					log.Println("pg service status error: ", err.Error()) // console log

					select {
					case <-ctx.Done():
						logger.Error("pg service is not ready: ", ctx.Err())
						log.Println("pg service is not ready: ", ctx.Err()) // console log
						return errors.Wrap(ctx.Err(), "pg service is not ready: "+err.Error())
					default:
						logger.Warn("retry init pg service")
						time.Sleep(5 * time.Second)
					}
				}
			},

			OnStop: func(ctx context.Context) error {
				return sqlService.Close()
			},
		})

		return sqlService, nil
	})
}
