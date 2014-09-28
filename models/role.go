package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

var (
	ROLE_TYPE_ALLOW  = 0
	ROLE_TYPE_FORBID = 1
)

var RoleType = map[int]string{
	ROLE_TYPE_ALLOW:  "允许",
	ROLE_TYPE_FORBID: "禁止",
}

type Role struct {
	Id        int64     `json:"id"`
	Name      string    `sql:"size:45" json:"name"`
	Auths     string    `sql:"size:20000" json:"auths"`
	Type      int       `json:"type"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	EditAt    time.Time `json:"edit_at"`
	TypeName  string    `sql:"-" json:"type_name"`
}

func getRoleDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func (r *Role) Refresh() error {
	db := getRoleDB()
	return db.First(r, r.Id).Error
}

func (r *Role) Save() error {
	db := getRoleDB()
	if r.Id == 0 {
		r.CreatedAt = time.Now()
	}
	r.EditAt = time.Now()
	return db.Save(r).Error
}

func (r *Role) Delete() error {
	db := getRoleDB()
	return db.Delete(&r).Error
}

func GetRoleList() (*[]Role, error) {
	var roles []Role
	db := getRoleDB()

	err := db.Find(&roles).Error
	if err == nil {
		for i := 0; i < len(roles); i++ {
			roles[0].TypeName = RoleType[roles[0].Type]
		}
	}

	return &roles, err
}

func GetRoleTypes() interface{} {
	var (
		kvs []TextValue
		kv  TextValue
	)
	for v, t := range RoleType {
		kv = TextValue{
			Text:  t,
			Value: v,
		}
		kvs = append(kvs, kv)
	}
	return &kvs
}
