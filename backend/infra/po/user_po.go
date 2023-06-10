package po

import (
	"fmt"
	"time"

	"code-shooting/infra/database/pg/sql"

	"github.com/pkg/errors"
)

type UserPo struct {
	Id         string `gorm:"primary_key" json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
	OrgId      string `json:"orgID"`
	TeamName   string `json:"teamName"`
	CenterName string `json:"centerName"`
	Institute  string `json:"institute"`
	Email      string `json:"email"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

type UserDB struct {
	*sql.GormDB
}

func (m UserDB) InsertUser(user *UserPo) error {
	//result := m.Exec("delete from mbr_User_po where name=? and ns=?", userName, userNs)
	result := m.Unscoped().Where("id=?", user.Id).Delete(user)
	if result.Error == nil {
		result = m.Model(user).Save(user)
	}
	return errors.WithStack(result.Error)
}

func (m UserDB) UpdateUser(user *UserPo) error {
	result := m.Model(user).Where("id=?", user.Id).Updates(user)
	return errors.WithStack(result.Error)
}

func (m UserDB) DeleteUser(user *UserPo) error {
	result := m.Model(user).Where("id=?", user.Id).Delete(user)
	return errors.WithStack(result.Error)
}

func (m UserDB) FindUsers(where, value string) ([]UserPo, error) {
	users := make([]UserPo, 0)
	result := m.Find(&users, fmt.Sprintf("%s=?", where), value)
	return users, errors.WithStack(result.Error)
}

func (m UserDB) FindUser(value string) (UserPo, error) {
	var u UserPo
	result := m.Where("id = ?", value).First(&u)
	return u, result.Error
}

func (m UserDB) IsExist(user *UserPo) (bool, error) {
	users := make([]UserPo, 0)
	result := m.Model(user).Where("id=?", user.Id).Find(&users)

	//result := m.Model(&UserPo{}).Where("User_json=?", userUserJson,).Find(&Users)
	if result.RowsAffected == 0 && result.Error == nil {
		return false, nil
	}
	if result.Error != nil {
		return false, errors.WithStack(result.Error)
	}
	return true, nil
}

func (m UserDB) FindAll() ([]UserPo, error) {
	users := make([]UserPo, 0)
	result := m.Model(&UserPo{}).Find(&users)
	return users, result.Error
}
