package userservice

import (
	"code-shooting/app/privilegeapp"
	"code-shooting/domain/entity"
	"code-shooting/domain/service/user"
	"code-shooting/infra/common"
	"code-shooting/infra/errcode"
	"code-shooting/infra/handler"
	"code-shooting/interface/dto"
	"fmt"
)

type GetInfo struct {
}

func NewGetInfoService() *GetInfo {
	return &GetInfo{}
}
func (g *GetInfo) GetUserInfoAndUpdate(request *entity.UserEntity) *common.Response {
	userInfo, err := user.NewUserDomainService().QueryUser(request)
	if err != nil {
		if errcode.SameCause(err, errcode.ErrRecordNotFound) {
			return hanldUser(request, "", func(userE *entity.UserEntity) error {
				return user.NewUserDomainService().AddUser(userE)
			})
		}
		return handler.Failure(err.Error())
	} else {
		userInfo.Id = request.Id
		updateUserWhenDiffProj(userInfo, userInfo.Department)
	}
	userInfo.Privileges = privilegeapp.GetStaffPrivileges(request.Id)
	return handler.SuccessGet(userInfo)
}

func (g *GetInfo) GetPerson(id string) (*dto.Person, error) {
	res := g.GetUserInfo(&entity.UserEntity{Id: id})
	if res.Result == handler.SUCCESS {
		user := ((*res).Detail).(*entity.UserEntity)
		personObj := dto.NewPerson(user)
		return &personObj, nil
	}
	return nil, fmt.Errorf("%s cant find user", id)
}

func (g *GetInfo) RefreshPersonMessage(id string) (*dto.Person, error) {
	userInfo, err := g.GetPerson(id)
	if err != nil {
		return nil, fmt.Errorf("%s cant find user", id)
	}
	res, err := getStaffFieldMessage(id)
	if err != nil {
		return nil, fmt.Errorf("%s cant find user", id)
	}
	userInfo.TeamName = res.TeamName
	userInfo.Department = res.Department
	userInfo.CenterName = res.CenterName
	userInfo.Institute = res.Institute
	return userInfo, nil
}

func (g *GetInfo) GetUserInfo(request *entity.UserEntity) *common.Response {
	userInfo, err := user.NewUserDomainService().QueryUser(request)
	if err != nil {
		if errcode.SameCause(err, errcode.ErrRecordNotFound) {
			return hanldUser(request, "", func(userE *entity.UserEntity) error {
				return user.NewUserDomainService().AddUser(userE)
			})
		}
		return handler.Failure(err.Error())
	}
	userInfo.Privileges = privilegeapp.GetStaffPrivileges(request.Id)
	return handler.SuccessGet(userInfo)
}
func updateUserWhenDiffProj(oldU *entity.UserEntity, projName string) {
	hanldUser(oldU, projName, func(userE *entity.UserEntity) error {
		return user.NewUserDomainService().ModifyUser(userE)
	})
}

type hanlduserInDb func(user *entity.UserEntity) error

func hanldUser(userInfo *entity.UserEntity, projName string, hanld hanlduserInDb) *common.Response {
	if len(projName) == 0 {
		var err error
		userInfo, err = getStaffFieldMessage(userInfo.Id)
		if err != nil {
			return handler.Failure(err.Error())
		}
		err1 := hanld(userInfo)
		if err1 != nil {
			return handler.Failure(err1.Error())
		}
	}
	userInfo.Privileges = privilegeapp.GetStaffPrivileges(userInfo.Id)
	return handler.SuccessCreated(userInfo)
}

func getStaffFieldMessage(id string) (*entity.UserEntity, error) {
	var person entity.UserEntity

	return &person, nil
}
