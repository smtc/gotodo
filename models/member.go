package models

import "github.com/jinzhu/gorm"

// 账号管理

type Member struct {
	Id         int64  `json:"id"`
	ObjectId   string `sql:"size:64" json:"object_id"`
	Name       string `sql:"size:40" json:"name"`
	Email      string `sql:"size:100" json:"email"`
	Avatar     string `sql:"size:120" json:"avatar"`
	Msisdn     string `sql:"size:20" json:"msisdn"`
	Password   string `sql:"size:80" json:"password"`
	Roles      string `sql:"type:text" json:"roles"` // 这是一个string数组, 以,分割
	Approved   bool   `json:"approved"`
	Activing   bool   `json:"acitiving"`
	ApprovedBy string `sql:"size:20" json:"approved_by"`
	IpAddr     string `sql:"size:30" json:"ipaddr"`
	DaysLogin  int    `json:"days_login"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	LastLogin int64 `json:"last_login"`

	Notifications int `json:"notifications"`
	Messages      int `json:"messages"`
}

func getMemberDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func (m *Member) Get(id int64) error {
	db := getMemberDB()
	return db.First(m, id).Error
}

func (m *Member) Save() error {
	db := getMemberDB()
	return db.Save(m).Error
}

func (m *Member) Delete() error {
	db := getMemberDB()
	return db.Delete(m).Error
}

func MemberDelete(where string) {
	db := getMemberDB()
	db.Where(where).Delete(&Member{})
}

func MemberList(page, size int, filter *map[string]interface{}) ([]Member, error) {
	db := getMemberDB()
	var accts []Member

	err := db.Offset(page * size).Limit(size).Find(&accts).Error
	return accts, err
}
