package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"code-shooting/infra/config"

	"code-shooting/infra/database/pg/sql"
	"code-shooting/infra/logger"
	"code-shooting/infra/restserver"

	"go.uber.org/fx"

	"code-shooting/app/privilegeapp"
	"code-shooting/app/service/result"
	"code-shooting/app/service/template"
	"code-shooting/domain"
	customuser "code-shooting/domain/service/custom-user"
	"code-shooting/domain/service/project"
	"code-shooting/domain/service/score"
	"code-shooting/infra"
	"code-shooting/infra/metricaop"
	"code-shooting/infra/privilegecfg"
	shootingresult "code-shooting/infra/shooting-result"
	"code-shooting/infra/shooting-result/defect"
	"code-shooting/infra/util"
	"code-shooting/infra/util/database"
	"code-shooting/interface/controller"
	"code-shooting/router"
)

var configPath = "/app/conf/code-shooting"

func init() {
	flag.StringVar(&configPath, "config-path", configPath, "code-shooting config file path")
	flag.StringVar(&util.DataDir, "data-dir", util.DataDir, "code shooting data dir")
	flag.StringVar(&util.ConfDir, "conf-dir", util.ConfDir, "code shooting config dir")
	flag.Parse()
}

func main() {
	if len(configPath) == 0 {
		log.Fatal("config path is empty")
		return
	}
	loggerInstance, err := logger.NewLogger(configPath)

	if err != nil {
		log.Fatal("failed create code-shooting logger: ", err)
		return
	}
	logger.SetLogger(loggerInstance.Named("code-shooting"))

	app := fx.New(
		fx.Logger(logger.GetLogger().CreateStdLogger()),
		fx.StartTimeout(15*time.Second),
		fx.StopTimeout(15*time.Second),
		config.NewModule(configPath),
		restserver.NewModule(),

		sql.NewModule(), fx.Invoke(database.InitGormDb),

		fx.Invoke(func() error {
			return privilegeapp.LoadPrivilegeCfg(filepath.Join(util.ConfDir, "privilege"))
		}),
		fx.Invoke(func(lc fx.Lifecycle) error {
			return privilegecfg.InvokeConfigWatch(lc, filepath.Join(util.ConfDir, "privilege"),
				privilegeapp.LoadPrivilegeCfg, project.ReloadProjectService)
		}),
		fx.Invoke(func(lc fx.Lifecycle) error {
			return privilegecfg.InvokeConfigWatch(lc, filepath.Join(util.ConfDir, "project"),
				customuser.ReloadCustomUserService)
		}),

		fx.Invoke(router.Register),

		fx.Invoke(func() {
			score.SetScoreService(shootingresult.NewShootingResultCalculator())
			controller.SetSubmitController(defect.NewDefectEncoder)
			result.SetResultService(func(workspace string, templateVersion string) (score.DefectCoder, error) {
				fileName, _ := controller.GetCurrentTemplateFileName(workspace, templateVersion)
				return defect.NewDefectEncoder(filepath.Join(util.TemplateDir, workspace, fileName))
			})
		}),

		domain.NewDomain(),

		infra.NewInfra(),

		metricaop.NewExporterModule(),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), app.StartTimeout())
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal("start app: ", err)
		return
	}

	logger.Info("code-shooting is started!")
	template.GetTemplateAppService().InitDefaultTemplate()
	go metricaop.SyncTargetDefectStat()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	stopCtx, cancel := context.WithTimeout(context.Background(), app.StopTimeout())
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Fatal("stop app: ", err)
		return
	}
}
