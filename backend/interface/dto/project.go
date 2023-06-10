package dto

type Project struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ProjectsRsp struct {
	Projects []Project `json:"projects"`
}
