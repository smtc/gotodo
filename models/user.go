package models

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/gocache"
	"github.com/smtc/goutils"
)

var (
	USER_CACHE_KEY = "gotodo_user"
)

type User struct {
	Id       int64  `json:"id"`
	ObjectId string `sql:"size:64" json:"object_id"`
	Name     string `sql:"size:40" json:"name"`
	Email    string `sql:"size:100" json:"email"`
	Avatar   string `sql:"size:120" json:"avatar"`
	Password string `sql:"size:80" json:"password"`
	Activing bool   `json:"activing"`
	IpAddr   string `sql:"size:30" json:"ip_addr"`
	Level    int    `json:"level"`
	Des      string `sql:"size:500" json:"des"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	LastLogin int64 `json:"last_login"`

	Role     string  `sql:"-" json:"role"`
	projects []int64 `sql:"-"`
}

func md5Encode(p string) string {
	p += "I2t"
	return fmt.Sprintf("%x", md5.Sum([]byte(p)))
}

func getUserCacheKey(id int64) string {
	k := strconv.FormatInt(id, 36)
	return fmt.Sprintf("%s_%s", USER_CACHE_KEY, k)
}

func getUserDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func getAllUserIds() ([]int64, bool) {
	var (
		cache = gocache.GetCache()
		suc   bool
		v     interface{}
		ids   []int64
	)
	v, suc = cache.Get(USER_CACHE_KEY)
	if !suc {
		return ids, false
	}

	ids, suc = v.([]int64)
	return ids, suc
}

// users 缓存id数组，实体根据id从单独的缓存中取
// 因为实际业务中，取单个user的次数远远多于去列表
func GetAllUsers() ([]User, error) {
	var (
		cache = gocache.GetCache()
		suc   bool
		err   error
		ids   []int64
		user  *User
		users []User
	)

	ids, suc = getAllUserIds()
	if suc {
		for _, id := range ids {
			user, err = GetUser(id)
			if err == nil {
				users = append(users, *user)
			}
		}
		return users, nil
	}

	// 从数据库中获取数据，第一次获取数据，缓存所有数据
	var (
		db = getUserDB()
		k  string
	)
	err = db.Find(&users).Order("level").Order("email").Error
	if err == nil {
		for i := 0; i < len(users); i++ {
			user = &users[i]
			user.setText()
			ids = append(ids, user.Id)
			k = getUserCacheKey(user.Id)
			cache.Store(k, user, 0)
		}
		cache.Set(USER_CACHE_KEY, ids, 0)
	}

	return users, err
}

func GetUser(id int64) (*User, error) {
	var (
		k     = getUserCacheKey(id)
		cache = gocache.GetCache()
		user  User
		suc   bool
		err   error
	)
	suc, err = cache.Retrieve(k, &user)
	if suc == false {
		db := getUserDB()
		err = db.First(&user, id).Error
		user.setText()
		if err != nil {
			cache.Store(k, &user, 0)
		}
	}

	return &user, err
}

func GetUserName(id int64) string {
	user, err := GetUser(id)
	if err == nil {
		return user.Name
	}
	return ""
}

func GetMultUserName(ids string) string {
	var (
		names []string
		id    int64
	)
	for _, s := range strings.Split(ids, ",") {
		id = goutils.ToInt64(s, 0)
		if id > 0 {
			names = append(names, GetUserName(id))
		}
	}
	return strings.Join(names, ",")
}

func (u *User) setText() {
	u.Role = ROLES[u.Level]
	u.Password = ""
}

func (u *User) Save() error {
	var (
		db    = getUserDB()
		err   error
		cache = gocache.GetCache()
		isNew = false
		old   User
	)

	if u.Id == 0 {
		isNew = true
		u.ObjectId = goutils.ObjectId()
		u.CreatedAt = time.Now().Unix()
		u.Password = md5Encode(u.Password)
	} else {
		err = db.First(&old, u.Id).Error
		if err != nil {
			return err
		}
		u.projects = old.projects

		// 如果未填写密码，使用旧密码
		if u.Password == "" {
			u.Password = old.Password
		} else {
			u.Password = md5Encode(u.Password)
		}
	}
	u.UpdatedAt = time.Now().Unix()

	err = db.Save(u).Error
	if err != nil {
		return err
	}
	u.setText()

	// 更新缓存
	k := getUserCacheKey(u.Id)
	cache.Store(k, u, 0)

	// 新建用户插入缓存id
	if isNew {
		ids, suc := getAllUserIds()
		if suc {
			ids = append(ids, u.Id)
			cache.Set(USER_CACHE_KEY, ids, 0)
		}
	}

	return nil
}

func UserDelete(id int64) (*User, error) {
	var (
		cache = gocache.GetCache()
		db    = getUserDB()
		err   error
		count int
	)

	// 判断用户是否执行过任务，如果执行过，不能删除
	db.Model(Task{}).Where("user=?", id).Count(&count)
	if count > 0 {
		user, err := GetUser(id)
		if err != nil {
			return user, err
		}
		user.Activing = false
		user.Save()
		return user, fmt.Errorf("用户已执行过任务，不能删除，已停用账号。")
	}

	err = db.Where("id=?", id).Delete(User{}).Error
	if err != nil {
		return nil, err
	}

	// 删除缓存
	k := getUserCacheKey(id)
	cache.Delete(k)

	// 从列表中移除
	ids, suc := getAllUserIds()
	if suc {
		var newids []int64
		for i, id := range ids {
			if id == id {
				newids = append(newids, ids[i+1:]...)
				break
			} else {
				newids = append(newids, id)
			}
		}

		cache.Set(USER_CACHE_KEY, newids, 0)
	}

	return nil, nil
}

func GetUserSelectData() interface{} {
	var (
		kvs   []TextValue
		kv    TextValue
		users []User
	)
	users, _ = GetAllUsers()
	for _, u := range users {
		kv = TextValue{
			Text:  u.Name,
			Value: u.Id,
		}
		kvs = append(kvs, kv)
	}
	return &kvs

}

func UserLogin(email, password string, r *http.Request) (*User, bool) {
	var (
		db   = getUserDB()
		err  error
		user User
	)

	password = md5Encode(password)
	err = db.Where("email=? and password=? and activing=1", email, password).First(&user).Error
	if err != nil {
		return nil, false
	}

	user.IpAddr = r.RemoteAddr
	user.LastLogin = time.Now().Unix()
	user.Password = ""
	user.Save()

	return &user, true
}
