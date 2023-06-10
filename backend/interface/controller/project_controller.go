package controller

import (
	"net/http"

	"code-shooting/infra/restserver"

	"code-shooting/domain/service/project"
	"code-shooting/interface/dto"
)

type ProjectController struct{}

func NewProjectController() *ProjectController {
	return &ProjectController{}
}

func (s *ProjectController) FindDeptsByProjectId(ctx restserver.Context) error {
	projectId := ctx.Param("id")
	depts := project.GetProjectService().FindDeptsByProjectId(projectId)
	return ctx.JSON(http.StatusOK, depts)
}

func (s *ProjectController) List(ctx restserver.Context) error {
	projs := project.GetProjectService().ListProjects()
	rsp := &dto.ProjectsRsp{Projects: make([]dto.Project, 0, len(projs))}
	for _, p := range projs {
		rsp.Projects = append(rsp.Projects, dto.Project{Id: p.Id, Name: p.Name})
	}
	return ctx.JSON(http.StatusOK, rsp)
}

func (s *ProjectController) FindProjectByUser(ctx restserver.Context) error {
	id := ctx.QueryParams().Get("userId")
	projs := project.GetProjectService().FindByUserId(id)
	rsp := &dto.ProjectsRsp{Projects: make([]dto.Project, 0, len(projs))}
	for _, p := range projs {
		rsp.Projects = append(rsp.Projects, dto.Project{Id: p.Id, Name: p.Name})
	}
	return ctx.JSON(http.StatusOK, rsp)
}
