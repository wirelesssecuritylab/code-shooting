package database

import (
	"code-shooting/infra/database/pg/sql"
	"code-shooting/infra/logger"

	"code-shooting/infra/po"
	metricpo "code-shooting/infra/po/metric-data-po"
)

var DB *sql.GormDB

func InitGormDb(sqlService sql.SqlService) {

	gormDB, err := sqlService.GetGormDB("cs_context")
	if err != nil {
		logger.Fatal("db abnormal", err)
		return
	}
	//defer gormDB.Close()

	err = gormDB.Migrator().AutoMigrate(
		new(po.UserPo), new(po.TargetPo), new(po.TbResult), new(po.TbRange), new(po.TemplatePo), new(po.TemplateOpHistoryPo),
		new(po.ShootingNotePo), new(po.ShootingDraftPo), new(metricpo.RangeRequestPo), new(metricpo.ShootingAccuracyPo),
		new(metricpo.ShootingRecordPo), new(metricpo.RingNumPo), new(metricpo.ShootingDurationPo), new(metricpo.TargetDefectStatPo),
		new(po.DefectPo),
		&po.TbEC{},
	)
	if err != nil {
		logger.Errorf("db auto migrate failed: %v", err)
		gormDB.Close()
		return
	}

	DB = gormDB
}
