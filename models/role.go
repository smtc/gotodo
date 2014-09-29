package models

import (
	"fmt"
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
	Id        int64  `json:"id"`
	Name      string `sql:"size:45" json:"name"`
	Auths     string `sql:"size:20000" json:"auths"`
	Type      int    `json:"type"`
	Des       string `sql:"size:512" json:"des"`
	CreatedBy int64  `json:"created_by"`
	CreatedAt int64  `json:"created_at" format:"datetime"`
	EditAt    int64  `json:"edit_at" format:"datetime"`
	TypeText  string `sql:"-" json:"type_text"`
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
		r.CreatedAt = time.Now().Unix()
	} else {
		old := Role{}
		db.First(&old, r.Id)
		//r.CreatedAt = old.CreatedAt
	}
	r.EditAt = time.Now().Unix()
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
			roles[0].TypeText = RoleType[roles[0].Type]
		}
	}

	return &roles, err
}

func RoleDelete(ids []int64) error {
	var (
		db    = getRoleDB()
		count = len(ids)
	)

	if count == 0 {
		return fmt.Errorf("id can't be null.")
	}

	return db.Where("id in (?)", ids).Delete(&Role{}).Error
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
