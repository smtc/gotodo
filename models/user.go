package models

import "github.com/jinzhu/gorm"

// 账号管理

type User struct {
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

func getUserDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func (u *User) Get(id int64) error {
	db := getUserDB()
	return db.First(u, id).Error
}

func (u *User) Save() error {
	db := getUserDB()
	return db.Save(u).Error
}

func (u *User) Delete() error {
	db := getUserDB()
	return db.Delete(u).Error
}

func UserDelete(where string) {
	db := getUserDB()
	db.Where(where).Delete(&User{})
}

func UserList(page, size int, filter *map[string]interface{}) ([]User, error) {
	db := getUserDB()
	var users []User

	err := db.Offset(page * size).Limit(size).Find(&users).Error
	return users, err
}
