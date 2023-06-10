package repository

import "code-shooting/domain/entity"

type UserInterface interface {
	InsertUser(entity *entity.UserEntity) error
	UpdateUser(entity *entity.UserEntity) error
	DeleteUser(entity *entity.UserEntity) error
	FindUsers(where, value string) ([]entity.UserEntity, error)
	FindUser(entity *entity.UserEntity) (*entity.UserEntity, error)
	IsExist(entity *entity.UserEntity) (bool, error)
	FindAll() ([]entity.UserEntity, error)
}
