package router

import (
	rs "code-shooting/infra/restserver"

	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"

	"code-shooting/infra/metricaop"
	"code-shooting/interface/controller"
)

func Register(multServer *rs.MultRestServer, mc *metricaop.MetricsCollector) error {
	server, err := multServer.GetRestServerByName("code-shooting")
	if err != nil {
		return errors.Wrap(err, "get server")
	}

	rg := server.RootGroupBox.RootGroup
	if rg == nil {
		return errors.New("without root group of server")
	}
	rg.POST("/actions/verify", controller.NewVerifyController().Verify)
	rg.POST("/actions/getInfo", controller.NewUserController().GetUserInfo)
	rg.POST("/actions/person", controller.NewUserController().Post)
	rg.POST("/answers/submit", controller.GetSubmitController().Submit, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(400)))
	rg.POST("/answers/save", controller.GetSubmitController().Save, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(400)))
	rg.GET("/answers/load", controller.GetSubmitController().Load, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(400)))
	rg.POST("/answers/savedraft", controller.GetSubmitController().SaveDraft, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(400)))
	rg.GET("/answers/loaddraft", controller.GetSubmitController().LoadDraft, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(400)))
	rg.GET("/answers/range/:id/language/:language", controller.NewRangeController().GetRangeAnswers)
	rg.GET("/results/range/:id/language/:language", controller.GetResultController().Get)
	rg.GET("/results/range/:id/language/:language/excel", controller.GetResultController().GetExcel)
	rg.GET("/projects", controller.NewProjectController().FindProjectByUser)
	rg.GET("/project/:id/depts", controller.NewProjectController().FindDeptsByProjectId)
	rg.POST("/actions/target", controller.NewTargetController().Post)
	rg.POST("/actions/target/:id/defect", controller.NewDefectController().Post)
	rg.GET("/actions/target/:workspace/defect/language", controller.NewDefectController().LoadSheetLangUpdate)
	rg.POST("/actions/target/:workspace/defect/language", controller.NewDefectController().GetLangDefect)
	rg.GET("/actions/target/:workspace", controller.NewDefectController().GetTemplateByWorksapce)
	rg.GET("/actions/target/:workspace/:template", controller.NewDefectController().GetLangsByWorkspaceAndTemplate)
	rg.POST("/actions/uploadTemp/owner/:id/type/:type", controller.NewTargetController().Upload)
	rg.POST("/actions/uploadTarget/target/:id/type/:type", controller.NewTargetController().UpdateFiles)
	rg.POST("/actions/range", controller.NewRangeController().Post)
	rg.GET("/range/shooted/user/:id/list", controller.NewRangeController().GetShootedRange)
	rg.POST("/targets/:id/files", controller.NewTargetController().GetCodeFile)
	rg.POST("/actions/target/exportbatchtarget", controller.NewTargetController().ExportBatchTarget)
	rg.GET("/targets/:id/answers/:rangeId", controller.NewTargetController().GetTargetAnswers)
	rg.GET("/documents/catalogue/:type/list", controller.NewDocumentController().GetDocmentsDir)
	rg.POST("/documents/detail/:catalogueId", controller.NewDocumentController().GetDocumentDetail)
	rg.POST("/actions/:workspace/uploadTemplate", controller.NewTemplateController().Upload)
	rg.POST("/actions/template", controller.NewTemplateController().Post)
	rg.GET("/workspaces", controller.NewWorkSpaceController().GetWorkSpaces)
	server.Use(mc.MetricsCollect)
	regECRouters(rg)
	return nil
}
