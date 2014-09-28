package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 账号管理

type Account struct {
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

	Birthday   time.Time `json:"birthday"`
	BannedAt   time.Time `json:"banned_at"`
	BannedTill time.Time `json:"banned_till"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ApprovedAt time.Time `json:"approved_at"`
	LastLogin  time.Time `json:"last_login"`
	LastPost   time.Time `json:"last_post"`

	Notifications int `json:"notifications"`
	Messages      int `json:"messages"`
}

func getAccountDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func (a *Account) Get(id int64) error {
	db := getAccountDB()
	return db.First(a, id).Error
}

func (a *Account) Save() error {
	db := getAccountDB()
	return db.Save(a).Error
}

func (a *Account) Delete() error {
	db := getAccountDB()
	return db.Delete(a).Error
}

func AccountDelete(where string) {
	db := getAccountDB()
	db.Where(where).Delete(&Account{})
}

func AccountList(page, size int, filter *map[string]interface{}) ([]Account, error) {
	db := getAccountDB()
	var accts []Account

	err := db.Offset(page * size).Limit(size).Find(&accts).Error
	return accts, err
}
