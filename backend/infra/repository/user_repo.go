package repository

import (
	"code-shooting/domain/entity"
	"code-shooting/infra/assembler"
	"code-shooting/infra/po"
	"code-shooting/infra/util/database"
)

type UserRepository struct {
	UserDB po.UserDB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		UserDB: po.UserDB{GormDB: database.DB},
	}
}

func (m *UserRepository) InsertUser(entity *entity.UserEntity) error {
	userPo := assembler.UserEntity2Po(entity)
	return m.UserDB.InsertUser(userPo)
}

func (m *UserRepository) UpdateUser(entity *entity.UserEntity) error {
	userPo := assembler.UserEntity2Po(entity)
	return m.UserDB.UpdateUser(userPo)
}

func (m *UserRepository) DeleteUser(entity *entity.UserEntity) error {
	userPo := assembler.UserEntity2Po(entity)
	return m.UserDB.DeleteUser(userPo)
}

func (m *UserRepository) FindUsers(where, value string) ([]entity.UserEntity, error) {
	Users, err := m.UserDB.FindUsers(where, value)
	return assembler.UserPos2Entities(Users), err
}

func (m *UserRepository) FindUser(value string) (*entity.UserEntity, error) {
	userPo, err := m.UserDB.FindUser(value)
	return assembler.UserPo2Entity(&userPo), cvtErr(err)
}

func (m *UserRepository) IsExist(entity *entity.UserEntity) (bool, error) {
	userPo := assembler.UserEntity2Po(entity)
	return m.UserDB.IsExist(userPo)
}

func (m *UserRepository) FindAll() ([]entity.UserEntity, error) {
	all, err := m.UserDB.FindAll()
	return assembler.UserPos2Entities(all), err
}
