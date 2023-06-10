package user

import (
	"fmt"

	"code-shooting/infra/logger"

	"code-shooting/domain/entity"
	"code-shooting/infra/errcode"
	"code-shooting/infra/repository"
)

type UserDomainService interface {
	AddUser(User *entity.UserEntity) error
	RemoveUser(User *entity.UserEntity) error
	ModifyUser(User *entity.UserEntity) error
	QueryUser(User *entity.UserEntity) (*entity.UserEntity, error)
	QueryAll() ([]entity.UserEntity, error)
	IsExistUser(User *entity.UserEntity) error
	IsNotExistUser(User *entity.UserEntity) error
	IsExist(User *entity.UserEntity) (bool, error)
}

type UserDomainServiceImpl struct {
	*repository.UserRepository
}

func NewUserDomainService() UserDomainService {
	return &UserDomainServiceImpl{
		repository.NewUserRepository(),
	}
}

func (s *UserDomainServiceImpl) AddUser(User *entity.UserEntity) error {
	if err := s.IsExistUser(User); err != nil {
		return err
	}
	if err := s.InsertUser(User); err != nil {
		logger.Errorf("AddUser %v failure cause:%s", User, err.Error())
		return err
	}

	logger.Infof("AddUser insert User %v", User)
	return nil
}

func (s *UserDomainServiceImpl) RemoveUser(User *entity.UserEntity) error {
	if err := s.IsNotExistUser(User); err != nil {
		logger.Errorf("RemoveUser ")
		return err
	}
	if err := s.DeleteUser(User); err != nil {
		logger.Errorf("RemoveUser %v failure cause:%s", User, err.Error())
		return err
	}
	return nil
}

func (s *UserDomainServiceImpl) ModifyUser(User *entity.UserEntity) error {
	if err := s.IsNotExistUser(User); err != nil {
		return err
	}
	if err := s.UpdateUser(User); err != nil {
		logger.Errorf("ModifyUser %v failure cause:%s", User, err.Error())
		return err
	}
	return nil
}

func (s *UserDomainServiceImpl) QueryAll() ([]entity.UserEntity, error) {
	return s.FindAll()
}

func (s *UserDomainServiceImpl) QueryUser(entity *entity.UserEntity) (*entity.UserEntity, error) {
	user, err := s.FindUser(entity.Id)
	if err != nil {
		logger.Warnf("QueryUser %v failure cause:%s", entity, err.Error())
		return nil, err
	}
	return user, nil
}

func (s *UserDomainServiceImpl) IsExistUser(User *entity.UserEntity) error {
	flag, err := s.IsExist(User)
	if err != nil {
		logger.Warnf("IsExistUser %s %v failure cause:%s", User, err.Error())
		return err
	}
	if flag {
		logger.Warnf("IsExistUser %s %v failure cause:%s", User, fmt.Sprintf("User:%s is exist", User.Id))
		return fmt.Errorf("user:%s already exist", User.Id)
	}
	return nil
}

func (s *UserDomainServiceImpl) IsNotExistUser(User *entity.UserEntity) error {
	flag, err := s.IsExist(User)
	if err != nil {
		logger.Warnf("IsNotExistUser %s %v failure cause:%s", User, err.Error())
		return err
	}
	if !flag {
		logger.Warnf("IsNotExistUser %s %v failure cause:%s", User, fmt.Sprintf("User:%s is not exist", User.Id))
		return errcode.ErrRecordNotFound
	}
	return nil
}
