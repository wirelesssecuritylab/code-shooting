package dto

import (
	"code-shooting/domain/entity"
)

type PersonAction struct {
	Action string `json:"name"`
	Params Person `json:"parameters"`
}

type Person struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
	OrgName    string `json:"orgName"`
	Institute  string `json:"institute"`
	Email      string `json:"email"`
	TeamName   string `json:"teamName"`
	CenterName string `json:"centerName"`
}

func NewPerson(u *entity.UserEntity) Person {
	if u == nil {
		return Person{}
	}
	projectName := ""
	return Person{Id: u.Id, Name: u.Name, Department: u.Department, OrgName: projectName, Institute: u.Institute, Email: u.Email, TeamName: u.TeamName, CenterName: u.CenterName}
}
