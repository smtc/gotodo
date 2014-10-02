package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/smtc/gocache"
	"github.com/smtc/goutils"
)

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

var USER_CACHE_KEY = "all_user_list"

func getUserCacheKey(id int64) string {
	k := strconv.FormatInt(id, 36)
	return fmt.Sprintf("gotodo_user_%s", k)
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
		user  User
		users []User
	)

	ids, suc = getAllUserIds()
	if suc {
		for _, id := range ids {
			user, err = GetUser(id)
			if err == nil {
				users = append(users, user)
			}
		}
		return users, nil
	}

	// 从数据库中获取数据，第一次获取数据，缓存所有数据
	var (
		db = getUserDB()
		k  string
	)
	err = db.Find(&users).Error
	if err == nil {
		for i := 0; i < len(users); i++ {
			user = users[i]
			ids = append(ids, user.Id)
			k = getUserCacheKey(user.Id)
			cache.Store(k, user, 0)
		}
		cache.Set(USER_CACHE_KEY, ids, 0)
	}

	return users, err
}

func GetUser(id int64) (User, error) {
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
		if err != nil {
			cache.Store(k, &user, 0)
		}
	}

	return user, err
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

func (u *User) Save() error {
	var (
		db    = getUserDB()
		err   error
		cache = gocache.GetCache()
		isNew = false
	)

	if u.Id == 0 {
		isNew = true
	}

	err = db.Save(u).Error
	if err != nil {
		return err
	}

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

func (u *User) Delete() error {
	var (
		cache = gocache.GetCache()
		db    = getUserDB()
		err   error
	)
	err = db.Delete(u).Error
	if err != nil {
		return err
	}

	// 删除缓存
	k := getUserCacheKey(u.Id)
	cache.Delete(k)

	// 从列表中移除
	ids, suc := getAllUserIds()
	if suc {
		var newids []int64
		for i, id := range ids {
			if id == u.Id {
				newids = append(newids, ids[i+1:]...)
				break
			} else {
				newids = append(newids, id)
			}
		}

		cache.Set(USER_CACHE_KEY, newids, 0)
	}

	return nil
}
